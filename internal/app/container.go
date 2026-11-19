package app

import (
	"context"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/adapters/postgres"
	"github.com/chekrzk/WhereToGoApp/internal/handlers/health"
	nearbyhandler "github.com/chekrzk/WhereToGoApp/internal/handlers/nearby"
	nearbyservice "github.com/chekrzk/WhereToGoApp/internal/services/nearby"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Pool           *pgxpool.Pool
	HealthHandler  *health.Handler
	FindPOIHandler *nearbyhandler.Handle
}

func (c *Container) Close() {
	c.Pool.Close()
}

func NewContainer(
	ctx context.Context,
	databaseURL string,
) (*Container, error) {
	pool, err := postgres.NewPool(ctx, databaseURL)

	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	repo := postgres.NewPOIRepository(pool)

	findNearbyService := nearbyservice.New(repo)

	healthHandler := health.New(pool)
	findNearbyHandler := nearbyhandler.New(findNearbyService)

	return &Container{
		Pool:           pool,
		HealthHandler:  healthHandler,
		FindPOIHandler: findNearbyHandler,
	}, nil
}
