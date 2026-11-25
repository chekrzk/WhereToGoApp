package buildinghandler

import (
	"context"

	"github.com/chekrzk/WhereToGoApp/internal/domain/building"
	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
)

type BuildSearcher interface {
	FindBuilding(context.Context, geo.Point) (building.ResolvedBuilding, error)
}
