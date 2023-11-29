package product

import (
	"context"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/uuid"
	"github.com/rs/zerolog"
)

type (
	Service struct {
		repository Repository
		logger     *zerolog.Logger
	}

	ServiceParams struct {
		Ctx        context.Context
		Repository Repository
		Logger     *zerolog.Logger
	}
)

func NewService(opt ServiceParams) *Service {
	return &Service{
		repository: opt.Repository,
		logger:     opt.Logger,
	}
}

func (s *Service) GetProductsByPredictionID(ctx context.Context, predictionID uuid.UUID) ([]Product, error) {
	return s.repository.GetProductsByPredictionID(ctx, predictionID)
}
