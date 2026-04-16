package random

import (
	"gitflic.ru/lms/internal/domain/quiz/criteria"
)

type RandomCriteria struct {
	count int
}

func NewRandomCriteria(params Params) (RandomCriteria, error) {
	if err := validateCount(params.Count); err != nil {
		return RandomCriteria{}, err
	}

	return RandomCriteria{
		count: params.Count,
	}, nil
}

func (c RandomCriteria) Count() int {
	return c.count
}

func (c RandomCriteria) QuestionCount() int {
	return c.Count()
}

func (c RandomCriteria) Type() criteria.CriteriaType {
	return criteria.CriteriaTypeRandom
}
