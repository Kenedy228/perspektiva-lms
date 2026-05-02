package matching

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/matching/pair"
	"gitflic.ru/lms/internal/domain/question/title"
)

type Question struct {
	*base.Base
	pairs []pair.Pair
}

func New(t title.Title, pairs []pair.Pair) (*Question, error) {
	base, err := base.New(t)
	if err != nil {
		return nil, err
	}

	if err := validatePairs(pairs); err != nil {
		return nil, err
	}

	return &Question{
		Base:  base,
		pairs: slices.Clone(pairs),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Pairs() []pair.Pair {
	return slices.Clone(q.pairs)
}

func (q *Question) Type() question.Type {
	return question.TypeMatching
}

func (q *Question) ChangePairs(pairs []pair.Pair) error {
	if err := validatePairs(pairs); err != nil {
		return err
	}

	q.pairs = slices.Clone(pairs)
	return nil
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:  q.Base.Clone(),
		pairs: q.Pairs(),
	}
}
