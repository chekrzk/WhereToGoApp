package app

import (
	"context"
	"fmt"
	"time"

	"github.com/chekrzk/WhereToGoApp/internal/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

type Options struct {
	Port        string
	DatabaseURL string
}

func Run(GracefulCtx context.Context, options Options, log zerolog.Logger) error {
	container, err := NewContainer(GracefulCtx, options.DatabaseURL)

	if err != nil {
		return fmt.Errorf("load container: %w", err)
	}
	defer container.Close()

	fiberApp := fiber.New()

	fiberApp.Use(middlewares.RequestID())
	fiberApp.Use(middlewares.AccessLogger(log))

	container.HealthHandler.RegisterRoutes(fiberApp)
	container.FindPOIHandler.RegisterRoutes(fiberApp)

	log.Info().
		Str("port", options.Port).
		Msg("Starting http server")

	if err := fiberApp.Listen(":"+options.Port,
		fiber.ListenConfig{
			GracefulContext: GracefulCtx,
			ShutdownTimeout: 10 * time.Second,
		},
	); err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", options.Port, err)
	}
	return nil
}
