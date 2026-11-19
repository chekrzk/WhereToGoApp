package poi

import (
	"errors"

	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
)

var ErrInvalidQuery = errors.New("invalid nearby query")

type Place struct {
	Location geo.Point

	ID string
	// H3ID string

	Name        string
	Category    string
	Subcategory string
}

type NearbyPlace struct {
	Place Place

	DistanceM float64
}

type NearbyQuery struct {
	Center  geo.Point
	RadiusM float64

	Limit       int
	Category    string
	Subcategory string
}
