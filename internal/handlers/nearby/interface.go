package nearbyhandler

import (
	"context"

	"github.com/chekrzk/WhereToGoApp/internal/domain/poi"
)

type POIsearcher interface {
	FindPOI(context.Context, poi.NearbyQuery) ([]poi.NearbyPlace, error)
}
