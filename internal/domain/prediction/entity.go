package prediction

import (
	"github.com/Gabukuro/insta-gift-api/internal/base"
	"github.com/uptrace/bun"
)

type (
	Prediction struct {
		bun.BaseModel   `bun:"predictions,alias:predictions"`
		base.BaseEntity `bun:"embed" json:"inline"`
		Username        string  `bun:"username" json:"username"`
		FeedbackRating  float64 `bun:"feedback_rating" json:"feedback_rating"`
	}
)
