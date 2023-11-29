package product

import (
	"github.com/Gabukuro/insta-gift-api/internal/base"
	"github.com/uptrace/bun"
)

type (
	Product struct {
		bun.BaseModel   `bun:"products,alias:products"`
		base.BaseEntity `bun:"embed" json:",inline"`
		Name            string  `bun:"name" json:"name"`
		Category        string  `bun:"main_category" json:"category"`
		Image           string  `bun:"image" json:"image"`
		Rating          float64 `bun:"ratings" json:"rating"`
		NoOfReviews     float64 `bun:"no_of_ratings" json:"no_of_reviews"`
	}
)
