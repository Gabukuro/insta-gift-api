package base

import (
	"time"

	"github.com/Gabukuro/insta-gift-api/internal/pkg/uuid"
)

type (
	BaseEntity struct {
		ID        uuid.UUID `bun:"id,pk,type:uuid" json:"id"`
		CreatedAt time.Time `bun:"created_at,default:now()" json:"created_at"`
		UpdatedAt time.Time `bun:"updated_at,default:now()" json:"updated_at"`
	}
)

func (b *BaseEntity) SetID(id uuid.UUID) {
	b.ID = id
}

func (b *BaseEntity) NewUUID() {
	id := uuid.New()
	b.SetID(id)
}

func (b *BaseEntity) SetCreatedAt() {
	timeNow := time.Now()
	b.CreatedAt = timeNow
	b.UpdatedAt = timeNow
}

func (b *BaseEntity) SetUpdatedAt() {
	timeNow := time.Now()
	b.UpdatedAt = timeNow
}
