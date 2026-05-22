package answer

import (
	"fmt"

	"github.com/google/uuid"
)

type OptionID struct {
	id uuid.UUID
}

func NewOptionID(id uuid.UUID) (OptionID, error) {
	if id == uuid.Nil {
		return OptionID{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return OptionID{id: id}, nil
}

func (o OptionID) ID() uuid.UUID {
	return o.id
}

func (o OptionID) IsZero() bool {
	return o.id == uuid.Nil
}
