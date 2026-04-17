package base

import (
	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Params struct {
	Text        question.QText
	Description question.QDescription
	ImageID     uuid.UUID
}
