package short

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
)

type Checker struct {
	normalizers []normalizer
}

func New(opts ...Option) Checker {
	c := Checker{}

	for i := range opts {
		if opts[i] != nil {
			opts[i](&c)
		}
	}

	return c
}

func (c Checker) Check(q question2.Question, a question2.Answer) (score.Score, error) {
	castQ, ok := q.(*short.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	castA, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	variants := castQ.Variants()
	studentAnswer := applyNormalizers(castA.Value(), c.normalizers...)

	for i := range variants {
		nVariant := applyNormalizers(variants[i].Text().Value(), c.normalizers...)

		if nVariant == studentAnswer {
			s, _ := score.New(1)
			return s, nil
		}
	}

	s, _ := score.New(0)
	return s, nil
}

func (c Checker) Supports(t question2.Type) bool {
	return t == question2.TypeShort
}
