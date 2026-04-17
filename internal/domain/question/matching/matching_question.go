package matching

import (
	"slices"

	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const minPairs = 2
const maxPairs = 20

type MatchingQuestion struct {
	base.Base
	pairs   []Pair
	options []Option
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validatePairs(params.Pairs, params.PairsCount); err != nil {
		return nil, err
	}

	pairs, options, err := mapPairs(params.Pairs)
	if err != nil {
		return nil, err
	}

	return &MatchingQuestion{
		Base:    base,
		pairs:   pairs,
		options: options,
	}, nil
}

func (q *MatchingQuestion) Pairs() []Pair {
	return slices.Clone(q.pairs)
}

func (q *MatchingQuestion) UpdatePairs(rawPairs map[string]content.RichContent, pairsCount int) error {
	if err := validatePairs(rawPairs, pairsCount); err != nil {
		return err
	}

	pairs, options, err := mapPairs(rawPairs)
	if err != nil {
		return err
	}

	q.pairs = pairs
	q.options = options
	q.Touch()
	return nil
}

func (q *MatchingQuestion) Options() []Option {
	return slices.Clone(q.options)
}

func (q *MatchingQuestion) Type() question.Type {
	return question.TypeMatching
}
