package app

import (
	"context"
	"fmt"

	postgresrepo "github.com/chekrzk/WhereToGoApp/internal/adapters/postgres"
	buildinghandler "github.com/chekrzk/WhereToGoApp/internal/handlers/building"
	"github.com/chekrzk/WhereToGoApp/internal/handlers/health"
	nearbyhandler "github.com/chekrzk/WhereToGoApp/internal/handlers/nearby"
	postgrespool "github.com/chekrzk/WhereToGoApp/internal/platform/postgres"
	buildingcontain "github.com/chekrzk/WhereToGoApp/internal/services/BuildingContain"
	nearbyservice "github.com/chekrzk/WhereToGoApp/internal/services/nearby"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Pool                *pgxpool.Pool
	HealthHandler       *health.Handler
	FindPOIHandler      *nearbyhandler.Handle
	FindBuildingHandler *buildinghandler.Handler
}

func (c *Container) Close() {
	c.Pool.Close()
}

func NewContainer(
	ctx context.Context,
	databaseURL string,
) (*Container, error) {
	pool, err := postgrespool.NewPool(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	repoPOI := postgresrepo.NewPOIRepository(pool)
	repoBuilding := postgresrepo.NewBuildingRepo(pool)

	findNearbyService := nearbyservice.New(repoPOI)
	findBuildingService := buildingcontain.New(repoBuilding)

	healthHandler := health.New(pool)
	findNearbyHandler := nearbyhandler.New(findNearbyService)
	findBuildingHandler := buildinghandler.New(findBuildingService)

	return &Container{
		Pool:                pool,
		HealthHandler:       healthHandler,
		FindPOIHandler:      findNearbyHandler,
		FindBuildingHandler: findBuildingHandler,
	}, nil
}
