package answer

import (
	"fmt"

	"github.com/google/uuid"
)

type PromptID struct {
	id uuid.UUID
}

func NewPromptID(id uuid.UUID) (PromptID, error) {
	if id == uuid.Nil {
		return PromptID{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return PromptID{id: id}, nil
}

func (p PromptID) ID() uuid.UUID {
	return p.id
}
