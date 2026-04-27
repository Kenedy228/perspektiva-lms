package collection

import "github.com/google/uuid"

type Identifiable interface {
	ID() uuid.UUID
}
