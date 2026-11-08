package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    serverConfig
	Logger loggerConfig
}

type serverConfig struct {
	Port string
}

type loggerConfig struct {
	Level string
}

func (e *Config) LoadEnv() error {
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("Failed to load env file: %w", err)
	}

	port, exist := os.LookupEnv("PORT")
	if !exist || port == "" {
		return errors.New("port not defined")
	}
	e.Env.Port = port

	level, exist := os.LookupEnv("LOG_LEVEL")
	if !exist || level == "" {
		return errors.New("log level not defined")
	}
	e.Logger.Level = level

	return nil
}
