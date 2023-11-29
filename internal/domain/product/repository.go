package product

import (
	"context"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

type (
	Repository interface {
		GetProductsByPredictionID(ctx context.Context, predictionID uuid.UUID) ([]Product, error)
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

func (r *postgresRepository) GetProductsByPredictionID(ctx context.Context, predictionID uuid.UUID) (products []Product, err error) {
	query := r.db.NewSelect().
		Model((*Product)(nil)).
		Join("JOIN prediction_products pp ON pp.product_id = products.id").
		Where("pp.prediction_id = ?", predictionID)

	err = query.Scan(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
