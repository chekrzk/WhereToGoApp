package app

import (
	"context"
	"fmt"
	"time"

	"github.com/chekrzk/WhereToGoApp/internal/handlers/health"
	"github.com/chekrzk/WhereToGoApp/internal/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func Run(GracefulCtx context.Context, port string, log zerolog.Logger) error {
	app := fiber.New()

	app.Use(middlewares.RequestID())
	app.Use(middlewares.AccessLogger(log))

	healthHandler := health.New()
	healthHandler.RegisterRoutes(app)

	log.Info().
		Str("port", port).
		Msg("Started http server")

	if err := app.Listen(":"+port,
		fiber.ListenConfig{
			GracefulContext: GracefulCtx,
			ShutdownTimeout: 10 * time.Second,
		},
	); err != nil {
		return fmt.Errorf("failed to listen on port %v : %w", port, err)
	}
	return nil
}
