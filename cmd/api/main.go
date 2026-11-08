package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chekrzk/WhereToGoApp/internal/app"
	"github.com/chekrzk/WhereToGoApp/internal/config"
	"github.com/chekrzk/WhereToGoApp/internal/config/logger"
)

func main() {
	var cfg config.Config

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cfg.LoadEnv(); err != nil {
		log.Printf("failed to init config: %v", err)
		return
	}

	appLogger, err := logger.New(cfg.Logger.Level)

	if err != nil {
		log.Printf("failed to init logger: %v", err)
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(
		signals,
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer signal.Stop(signals)

	go func() {
		sig := <-signals

		appLogger.Info().
			Str("signal", sig.String()).
			Msg("shutdown signal received")

		cancel()
	}()

	if err := app.Run(ctx, cfg.Env.Port, appLogger); err != nil {
		appLogger.Error().
			Err(err).
			Msg("application stopped with error")
		return
	}
	appLogger.Info().
		Msg("application stopped gracefully")
}
