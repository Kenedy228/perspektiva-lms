package utils

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateID() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate id error, details: %w", err)
	}

	return id, nil
}
