package typed

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	qtyped "gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/question/typed/answer"
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
	qCast, ok := q.(*qtyped.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	blanks := qCast.Blanks()
	studentAnswers := aCast.BlanksAsMap()

	if len(studentAnswers) != len(blanks) {
		s, _ := score.New(0)
		return s, nil
	}

	for i := range blanks {
		studentVariant, ok := studentAnswers[blanks[i].Placeholder()]

		if !ok {
			s, _ := score.New(0)
			return s, nil
		}

		studentVariant = applyNormalizers(studentVariant, c.normalizers...)

		has := slices.ContainsFunc(blanks[i].VariantsValues(), func(current string) bool {
			current = applyNormalizers(current, c.normalizers...)
			return current == studentVariant
		})

		if !has {
			s, _ := score.New(0)
			return s, nil
		}
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question2.Type) bool {
	return qType == question2.TypeTyped
}
