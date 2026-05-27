package matching

import (
	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
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

	if _, ok := q.(*matching.Question); !ok {
		return grading.ErrInvalidQuestionType
	}

	if _, ok := a.(answer.Answer); !ok {
		return grading.ErrInvalidAnswerType
	}

	return nil
}
