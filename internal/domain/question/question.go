package question

import (
	"time"

	"github.com/google/uuid"
)

type Question interface {
	ID() uuid.UUID
	Text() string
	Instruction() string
	ImageID() uuid.UUID
	Type() Type
	CreatedAt() time.Time
	UpdatedAt() time.Time
	HasImage() bool
	Clone() Question
	CheckAnswer(Answer) bool
}
