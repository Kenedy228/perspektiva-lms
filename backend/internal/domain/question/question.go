package question

import (
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
)

type Question interface {
	ID() uuid.UUID
	Title() title.Title
	Instruction() string
	Type() Type
	Clone() Question
	ChangeTitle(title.Title) error
}
