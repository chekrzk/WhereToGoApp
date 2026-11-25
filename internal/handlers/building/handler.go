package buildinghandler

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/chekrzk/WhereToGoApp/internal/domain/building"
	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service BuildSearcher
}

func New(service BuildSearcher) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FindBuildingContain(ctx fiber.Ctx) error {

	var (
		place geo.Point
		err   error
	)

	place.Lat, err = strconv.ParseFloat(ctx.Query("lat"), 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "latitude has to be a number",
		})
	}

	place.Lon, err = strconv.ParseFloat(ctx.Query("lon"), 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "longitude has to be a number",
		})
	}

	searchCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	build, err := h.service.FindBuilding(searchCtx, place)
	if err != nil {
		if errors.Is(err, building.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if errors.Is(err, geo.ErrWrongCoordinate) {
			return ctx.Status(422).JSON(fiber.Map{
				"error": geo.ErrWrongCoordinate.Error(),
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to perform search",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(build)
}

func (h *Handler) RegisterRoute(router fiber.Router) {
	router.Get("/api/v1/buildings/resolve", h.FindBuildingContain)
}
