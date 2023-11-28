package prediction

import (
	"context"

	"github.com/rs/zerolog"
)

type (
	Service struct {
		repository Repository
		logger     *zerolog.Logger
	}
)

func NewService(ctx context.Context, repository Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error {
	return s.repository.GetPredictionByUsername(ctx, username, prediction)
}

func (s *Service) EnqueuePrediction(ctx context.Context, username string) error {
	return nil
}
