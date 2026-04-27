package matching

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
)

type Answer struct {
	pairs []Pair
}

func NewAnswer(params AnswerParams) question.Answer {
	cPairs := slices.Clone(params.Pairs)

	return Answer{
		pairs: cPairs,
	}
}

func (a Answer) Pairs() []Pair {
	return slices.Clone(a.pairs)
}

func (a Answer) IsEmpty() bool {
	return len(a.pairs) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		pairs: a.Pairs(),
	}
}
