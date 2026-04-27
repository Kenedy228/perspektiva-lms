package matching

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

type Question struct {
	*base.Base
	pairs []Pair
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validatePairs(params.Pairs); err != nil {
		return nil, err
	}

	cPairs := slices.Clone(params.Pairs)

	return &Question{
		Base:  base,
		pairs: cPairs,
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Pairs() []Pair {
	return slices.Clone(q.pairs)
}

func (q *Question) UpdatePairs(pairs []Pair) error {
	if err := validatePairs(pairs); err != nil {
		return err
	}

	cPairs := slices.Clone(pairs)
	q.pairs = cPairs
	q.Touch()
	return nil
}

func (q *Question) Type() question.Type {
	return question.TypeMatching
}

func (q *Question) ShuffledPairs() []Pair {
	cPairs := slices.Clone(q.pairs)
	rand.Shuffle(len(cPairs), func(i, j int) {
		cPairs[i], cPairs[j] = cPairs[j], cPairs[i]
	})

	return cPairs
}

func (q *Question) CheckAnswer(answer question.Answer) bool {
	castAns, ok := answer.(Answer)
	if !ok {
		return false
	}

	if len(q.pairs) != len(castAns.pairs) {
		return false
	}

	for i := range castAns.pairs {
		for j := i + 1; j < len(castAns.pairs); j++ {
			if castAns.pairs[i].Equal(castAns.pairs[j]) {
				return false
			}
		}
	}

	for i := range castAns.pairs {
		isPresent := slices.ContainsFunc(q.pairs, func(p Pair) bool {
			return p.Equal(castAns.pairs[i])
		})

		if !isPresent {
			return false
		}
	}

	return true
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:  q.Base.Clone(),
		pairs: q.Pairs(),
	}
}
