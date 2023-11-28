package prediction

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type (
	Repository interface {
		GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error
	}

	postgresRepository struct {
		logger *zerolog.Logger
		db     bun.DB
	}
)

func NewRepository(db *bun.DB, logger *zerolog.Logger) Repository {
	return &postgresRepository{
		logger: logger,
		db:     *db,
	}
}

func (r *postgresRepository) GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error {
	query := r.db.NewSelect().
		Model(prediction).
		Where("username = ?", username)

	return query.Scan(ctx)
}
