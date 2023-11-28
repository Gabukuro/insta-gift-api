package base

import (
	"time"

	"github.com/transfeera/go-infra/pkg/uuid"
)

type (
	BaseEntity struct {
		ID        uuid.UUID `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
		CreatedAt time.Time `bun:"created_at,default:now()" json:"created_at"`
		UpdatedAt time.Time `bun:"updated_at,default:now()" json:"updated_at"`
	}
)
