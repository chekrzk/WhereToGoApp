package buildingcontain

import (
	"context"
	"errors"
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

func (s *Service) FindBuilding(ctx context.Context, place geo.Point) (building.ResolvedBuilding, error) {

	if err := place.Validate(); err != nil {
		return building.ResolvedBuilding{}, fmt.Errorf("find building contain service: %w", err)
	}

	build, err := s.repo.FindContaining(ctx, place)

	if errors.Is(err, building.ErrNotFound) {
		return s.repo.FindNearest(ctx, place, 50)
	}
	if err != nil {
		return building.ResolvedBuilding{}, fmt.Errorf("find building contain service: %w", err)
	}

	return build, nil
}
