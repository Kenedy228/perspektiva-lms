package question

import (
	"time"

	"github.com/google/uuid"
)

type Question interface {
	ID() uuid.UUID
	Text() string
	Description() string
	Image() uuid.UUID
	HasImage() bool
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Type() Type
}
