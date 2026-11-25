package buildingcontain

import (
	"context"

	"github.com/chekrzk/WhereToGoApp/internal/domain/building"
	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
)

type BuildingFinder interface {
	FindContaining(context.Context, geo.Point) (building.ResolvedBuilding, error)
	FindNearest(context.Context, geo.Point, int) (building.ResolvedBuilding, error)
}
