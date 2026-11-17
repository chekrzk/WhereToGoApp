package health

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Pinger interface {
	Ping(ctx context.Context) error
}
type Handler struct {
	db Pinger
}

func New(db Pinger) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) Live(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{"status": "ok"},
	)
}

func (h *Handler) Ping(c fiber.Ctx) error {
	pingCtx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	err := h.db.Ping(pingCtx)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(
			fiber.Map{
				"status":   "not ready",
				"database": "unavailable",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status":   "ready",
			"database": "available",
		},
	)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/health/live", h.Live)
	router.Get("/health/ping", h.Ping)
}
