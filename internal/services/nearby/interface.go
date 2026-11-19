package nearbyservice

import (
	"context"

	"github.com/chekrzk/WhereToGoApp/internal/domain/poi"
)

type Finder interface {
	FindNearby(context.Context, poi.NearbyQuery) ([]poi.NearbyPlace, error)
}
