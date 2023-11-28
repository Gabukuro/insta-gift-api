package product

import (
	"github.com/Gabukuro/insta-gift-api/internal/base"
	"github.com/uptrace/bun"
)

type (
	Product struct {
		bun.BaseModel   `bun:"products,alias:products"`
		base.BaseEntity `bun:"embed" json:"inline"`
		Name            string `bun:"name" json:"name"`
		Category        string `bun:"category" json:"category"`
		Image           string `bun:"image" json:"image"`
		Rating          int    `bun:"rating" json:"rating"`
		NoOfReviews     int    `bun:"no_of_reviews" json:"no_of_reviews"`
	}
)
