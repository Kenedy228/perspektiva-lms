package question

import (
	"time"

	"github.com/google/uuid"
)

// NOTE: interface for Question polymorphism
type Question interface {
	ID() uuid.UUID
	Text() QText
	Description() QDescription
	ImageID() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Type() Type
	HasImage() bool
}
