package matching

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const minPairs = 2
const maxPairs = 20

type MatchingQuestion struct {
	base.Base
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

	pairs, err := mapPairs(params.Pairs)
	if err != nil {
		return nil, err
	}

	return &MatchingQuestion{
		Base:  base,
		pairs: pairs,
	}, nil
}

func (q *MatchingQuestion) Pairs() []Pair {
	return slices.Clone(q.pairs)
}

func (q *MatchingQuestion) UpdatePairs(rawPairs []PairParams) error {
	if err := validatePairs(rawPairs); err != nil {
		return err
	}

	pairs, err := mapPairs(rawPairs)
	if err != nil {
		return err
	}

	q.pairs = pairs
	q.Touch()
	return nil
}

func (q *MatchingQuestion) Type() question.Type {
	return question.TypeMatching
}

func (q *MatchingQuestion) HasPair(pair Pair) bool {
	return slices.Contains(q.pairs, pair)
}

func (q *MatchingQuestion) ShuffledPairs() []Pair {
	cPairs := slices.Clone(q.pairs)
	rand.Shuffle(len(cPairs), func(i, j int) {
		cPairs[i], cPairs[j] = cPairs[j], cPairs[i]
	})

	return cPairs
}
