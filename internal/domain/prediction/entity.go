package prediction

import (
	"github.com/Gabukuro/insta-gift-api/internal/base"
	"github.com/uptrace/bun"
)

type (
	Prediction struct {
		bun.BaseModel   `bun:"predictions,alias:predictions"`
		base.BaseEntity `bun:"embed" json:",inline"`
		Username        string           `bun:"username" json:"username"`
		FeedbackRating  int64            `bun:"feedback_rate" json:"feedback_rating"`
		Status          PredictionStatus `bun:"status" json:"status"`
	}

	PredictionStatus string
)

const (
	PredictionStatusPending   = PredictionStatus("pending")
	PredictionStatusCompleted = PredictionStatus("completed")
	PredictionStatusFailed    = PredictionStatus("failed")
)
