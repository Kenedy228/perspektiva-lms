package selectable

import (
	"gitflic.ru/lms/internal/domain/question"
)

type OptionParams struct {
	Content   question.Content
	IsCorrect bool
}
