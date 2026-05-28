package short

import (
	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
)

// Checker проверяет короткие текстовые ответы с опциональной нормализацией.
type Checker struct {
	normalizers []normalizer
}

// New создает checker коротких ответов и применяет переданные опции.
func New(opts ...Option) Checker {
	c := Checker{}

	for i := range opts {
		if opts[i] != nil {
			opts[i](&c)
		}
	}

	return c
}

// Check возвращает 1, если нормализованный ответ совпал с любым вариантом вопроса.
func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	castQ, ok := q.(*short.Question)
	if !ok {
		return score.Score{}, grading.ErrInvalidQuestionType
	}

	castA, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, grading.ErrInvalidAnswerType
	}

	variants := castQ.Variants()
	studentAnswer := applyNormalizers(castA.Value(), c.normalizers...)

	for i := range variants {
		nVariant := applyNormalizers(variants[i].Value(), c.normalizers...)
		if nVariant == studentAnswer {
			return score.New(1)
		}
	}

	return score.New(0)
}
