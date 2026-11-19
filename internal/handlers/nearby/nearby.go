package nearbyhandler

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/chekrzk/WhereToGoApp/internal/domain/poi"
	"github.com/gofiber/fiber/v3"
)

type Handle struct {
	service POIsearcher
}

func New(service POIsearcher) *Handle {
	return &Handle{
		service: service,
	}
}

func (h *Handle) FindNearby(ctx fiber.Ctx) error {
	var (
		query poi.NearbyQuery
		err   error
	)

	query.Center.Lat, err = strconv.ParseFloat(ctx.Query("lat"), 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "latitude has to be a number",
		})
	}

	query.Center.Lon, err = strconv.ParseFloat(ctx.Query("lon"), 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "longitude has to be a number",
		})
	}

	query.RadiusM, err = strconv.ParseFloat(ctx.Query("radius"), 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "radius has to be a number",
		})
	}

	query.Limit, err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit has to be a number",
		})
	}

	query.Category = ctx.Query("category")
	query.Subcategory = ctx.Query("subcategory")

	searchCtx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	places, err := h.service.FindPOI(searchCtx, query)
	if err != nil {
		if errors.Is(err, poi.ErrInvalidQuery) {
			return ctx.Status(422).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to perform search",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(places)
}

func (h *Handle) RegisterRoutes(router fiber.Router) {
	router.Get("/api/v1/pois/nearby/", h.FindNearby)
}
