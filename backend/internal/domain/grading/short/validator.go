package short

import (
	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
)

type Validator struct{}

func NewValidator() Validator {
	return Validator{}
}

func (v Validator) Validate(q question.Question, a question.Answer) error {
	if q == nil {
		return grading.ErrNilQuestion
	}

	if a == nil {
		return grading.ErrNilAnswer
	}

	if _, ok := q.(*short.Question); !ok {
		return grading.ErrInvalidQuestionType
	}

	if _, ok := a.(answer.Answer); !ok {
		return grading.ErrInvalidAnswerType
	}

	return nil
}
