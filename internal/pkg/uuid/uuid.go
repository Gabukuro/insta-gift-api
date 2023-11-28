package uuid

import (
	"github.com/gofrs/uuid"
)

type (
	UUID = uuid.UUID
)

func New() UUID {
	id, err := uuid.NewV6()
	if err != nil {
		panic(err)
	}

	return id
}

func Parse(s string) (UUID, error) {
	return uuid.FromString(s)
}

func MustParse(s string) UUID {
	id, err := uuid.FromString(s)
	if err != nil {
		panic(err)
	}

	return id
}

func IsNil(u UUID) bool {
	return u == uuid.Nil
}
