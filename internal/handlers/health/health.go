package health

import "github.com/gofiber/fiber/v3"

type Handler struct {
	//
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Live(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{"status": "ok"},
	)
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/health/live", h.Live)
}
