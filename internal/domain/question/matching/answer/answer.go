package answer

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Answer struct {
	pairs []AnswerPair
}

func New(pairs []AnswerPair) Answer {
	return Answer{
		pairs: slices.Clone(pairs),
	}
}

func (a Answer) Pairs() []AnswerPair {
	return slices.Clone(a.pairs)
}

func (a Answer) PairsAsMap() map[uuid.UUID]uuid.UUID {
	pairs := make(map[uuid.UUID]uuid.UUID, len(a.pairs))

	for i := range a.pairs {
		pairs[a.pairs[i].PromptID] = a.pairs[i].MatchID
	}

	return pairs
}

func (a Answer) IsEmpty() bool {
	return len(a.pairs) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		pairs: a.Pairs(),
	}
}
