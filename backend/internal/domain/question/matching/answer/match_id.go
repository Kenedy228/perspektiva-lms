package answer

import (
	"fmt"

	"github.com/google/uuid"
)

type MatchID struct {
	id uuid.UUID
}

func NewMatchID(id uuid.UUID) (MatchID, error) {
	if id == uuid.Nil {
		return MatchID{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return MatchID{id: id}, nil
}

func (m MatchID) ID() uuid.UUID {
	return m.id
}
