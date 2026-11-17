package app

import (
	"context"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/adapters/postgres"
	"github.com/chekrzk/WhereToGoApp/internal/handlers/health"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Pool          *pgxpool.Pool
	HealthHandler *health.Handler
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

	healthHandler := health.New(pool)

	return &Container{
		Pool:          pool,
		HealthHandler: healthHandler,
	}, nil
}
