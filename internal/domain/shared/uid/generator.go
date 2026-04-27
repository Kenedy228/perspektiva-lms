package uid

import (
	"github.com/google/uuid"
)

func New() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, ErrUIDGeneration
	}

	return id, nil
}
