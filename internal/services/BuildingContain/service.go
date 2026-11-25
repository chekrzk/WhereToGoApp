package buildingcontain

import (
	"context"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/domain/building"
	"github.com/chekrzk/WhereToGoApp/internal/domain/geo"
)

type Service struct {
	repo BuildingFinder
}

func New(repo BuildingFinder) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindBuilding(ctx context.Context, place geo.Point) (building.Building, error) {

	if err := place.Validate(); err != nil {
		return building.Building{}, fmt.Errorf("find building contain service: %w", err)
	}

	return s.repo.FindContaining(ctx, place)
}
