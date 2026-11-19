package nearbyservice

import (
	"context"
	"fmt"

	"github.com/chekrzk/WhereToGoApp/internal/domain/poi"
)

type Service struct {
	repo Finder
}

func New(repo Finder) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) FindPOI(ctx context.Context, query poi.NearbyQuery) ([]poi.NearbyPlace, error) {

	if err := query.Center.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", poi.ErrInvalidQuery, err)
	}

	if !(query.RadiusM >= 1 && query.RadiusM <= 5000) {
		return nil, fmt.Errorf("find poi service: %w", poi.ErrInvalidQuery)
	}

	if !(query.Limit >= 1 && query.Limit <= 100) {
		return nil, fmt.Errorf("find poi service: %w", poi.ErrInvalidQuery)
	}

	return s.repo.FindNearby(ctx, query)
}
