package prediction

import (
	"context"

	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type (
	Repository interface {
		CreatePrediction(ctx context.Context, prediction *Prediction) error
		GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error
	}

	postgresRepository struct {
		logger *zerolog.Logger
		db     bun.DB
	}
)

var pqErr *pq.Error

func NewRepository(db *bun.DB, logger *zerolog.Logger) Repository {
	repo := &postgresRepository{
		logger: logger,
		db:     *db,
	}

	return repo
}

func (r *postgresRepository) CreatePrediction(ctx context.Context, prediction *Prediction) error {
	_, err := r.db.NewInsert().Model(prediction).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresRepository) GetPredictionByUsername(ctx context.Context, username string, prediction *Prediction) error {
	query := r.db.NewSelect().
		Model(prediction).
		Where("username = ?", username)

	return query.Scan(ctx)
}
