package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rs/zerolog"
)

func AccessLogger(log zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()
		status := c.StatusCode()

		if err != nil {
			status = fiber.StatusInternalServerError

			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				status = fiberErr.Code
			}
		}

		event := log.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}

		if err != nil {
			event = event.Err(err)
		}

		event.
			Str("request_id", requestid.FromContext(c)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", status).
			Dur("duration", c.Elapsed()).
			Msg("HTTP request")

		return err
	}
}
