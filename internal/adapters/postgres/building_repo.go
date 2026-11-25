package postgresrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/domain/building"
	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BuildingRepo struct {
	pool *pgxpool.Pool
}

func NewBuildingRepo(pool *pgxpool.Pool) *BuildingRepo {
	return &BuildingRepo{
		pool: pool,
	}
}

func (r *BuildingRepo) FindContaining(ctx context.Context, place geo.Point) (building.Building, error) {

	var buildContain building.Building

	err := r.pool.QueryRow(ctx, `
	SELECT osm_type, osm_id, building_type, COALESCE(street, ''), COALESCE(house_number, '')
	FROM osm.buildings
	WHERE ST_Covers(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
	ORDER BY ST_Area(geom::geography), osm_type, osm_id
	LIMIT 1;
	`, place.Lon, place.Lat).Scan(&buildContain.Osm_type, &buildContain.Osm_id,
		&buildContain.Building_type, &buildContain.Street, &buildContain.House_number)

	if errors.Is(err, pgx.ErrNoRows) {
		return building.Building{}, fmt.Errorf("query find building: %w", building.ErrNotFound)
	}

	if err != nil {
		return building.Building{}, fmt.Errorf("query find building: %w", err)
	}

	return buildContain, nil
}
