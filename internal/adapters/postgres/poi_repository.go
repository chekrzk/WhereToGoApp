package postgres

import (
	"context"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
	"github.com/chekrzk/WhereToGoApp/internal/domain/poi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type POIRepository struct {
	pool *pgxpool.Pool
}

func NewPOIRepository(pool *pgxpool.Pool) *POIRepository {
	return &POIRepository{
		pool: pool,
	}
}

func (r *POIRepository) FindNearby(ctx context.Context, query poi.NearbyQuery) ([]poi.NearbyPlace, error) {
	const sql = `
		WITH search AS (
    		SELECT ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography AS center
		)
		SELECT
    		concat(p.osm_type, '/', p.osm_id) AS id,
    		COALESCE(p.name, '') AS name,
    		p.category,
    		p.subcategory,
    		ST_Y(p.geom) AS lat,
    		ST_X(p.geom) AS lon,
    		ST_Distance(p.geom::geography, search.center) AS distance_m
		FROM osm.pois AS p
		CROSS JOIN search
		WHERE ST_DWithin(p.geom::geography, search.center, $3)
 			AND ($4::text = '' OR p.category = $4)
			AND ($5::text = '' OR p.subcategory = $5)
		ORDER BY distance_m, p.osm_type, p.osm_id
		LIMIT $6;`

	rows, err := r.pool.Query(ctx, sql, query.Center.Lon, query.Center.Lat, query.RadiusM, query.Category, query.Subcategory, query.Limit)

	if err != nil {
		return nil, fmt.Errorf("querry find nearby: %w", err)
	}
	defer rows.Close()

	places := make([]poi.NearbyPlace, 0)

	for rows.Next() {
		var (
			id          string
			name        string
			category    string
			subcategory string
			lat         float64
			lon         float64
			distanceM   float64
		)

		if err := rows.Scan(&id, &name, &category, &subcategory, &lat, &lon, &distanceM); err != nil {
			return nil, fmt.Errorf("scan nearby places: %w", err)
		}

		places = append(places, poi.NearbyPlace{
			Place: poi.Place{
				Location: geo.Point{
					Lat: lat,
					Lon: lon,
				},
				ID:          id,
				Name:        name,
				Category:    category,
				Subcategory: subcategory,
			},
			DistanceM: distanceM,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read nearby places: %w", err)
	}

	return places, nil
}
