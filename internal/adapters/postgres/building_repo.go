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

func (r *BuildingRepo) FindContaining(ctx context.Context, place geo.Point) (building.ResolvedBuilding, error) {

	var buildContain building.ResolvedBuilding

	err := r.pool.QueryRow(ctx, `
	SELECT osm_type, osm_id, building_type, COALESCE(street, ''), COALESCE(house_number, '')
	FROM osm.buildings
	WHERE ST_Covers(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
	ORDER BY ST_Area(geom::geography), osm_type, osm_id
	LIMIT 1;
	`, place.Lon, place.Lat).Scan(&buildContain.Building.Osm_type, &buildContain.Building.Osm_id,
		&buildContain.Building.Building_type, &buildContain.Building.Street, &buildContain.Building.House_number)

	buildContain.DistanceM = 0
	buildContain.Match = "contains"

	if errors.Is(err, pgx.ErrNoRows) {
		return building.ResolvedBuilding{}, fmt.Errorf("query find building: %w", building.ErrNotFound)
	}

	if err != nil {
		return building.ResolvedBuilding{}, fmt.Errorf("query find building: %w", err)
	}

	return buildContain, nil
}

func (r *BuildingRepo) FindNearest(ctx context.Context, place geo.Point, radius int) (building.ResolvedBuilding, error) {
	var buildContain building.ResolvedBuilding

	err := r.pool.QueryRow(ctx, `
	WITH click AS (
    SELECT ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography AS point
	)
	SELECT
    	b.osm_type,
    	b.osm_id,
    	b.building_type,
    	COALESCE(b.street, ''),
    	COALESCE(b.house_number, ''),
    	ST_Distance(b.geom::geography, click.point) AS distance_m
	FROM osm.buildings AS b
	CROSS JOIN click
	WHERE ST_DWithin(b.geom::geography, click.point, $3)
	ORDER BY distance_m, b.osm_type, b.osm_id
	LIMIT 1;
	`, place.Lon, place.Lat, radius).Scan(&buildContain.Building.Osm_type,
		&buildContain.Building.Osm_id,
		&buildContain.Building.Building_type,
		&buildContain.Building.Street,
		&buildContain.Building.House_number,
		&buildContain.DistanceM)
	buildContain.Match = "nearest"

	if errors.Is(err, pgx.ErrNoRows) {
		return building.ResolvedBuilding{}, fmt.Errorf("query find nearest building: %w", building.ErrNotFound)
	}

	if err != nil {
		return building.ResolvedBuilding{}, fmt.Errorf("query find building: %w", err)
	}

	return buildContain, nil
}
