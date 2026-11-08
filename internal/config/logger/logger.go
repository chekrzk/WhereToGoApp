package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// TODO исправить, я не правильно сделал
func New(level string) (zerolog.Logger, error) {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	lvl, err := zerolog.ParseLevel(level)

	if err != nil {

		return zerolog.Logger{}, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	log := zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Logger()

	return log, nil
}
