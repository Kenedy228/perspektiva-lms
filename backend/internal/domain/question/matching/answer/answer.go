package answer

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
)

type Pair struct {
	PromptID PromptID
	MatchID  MatchID
}

type Answer struct {
	pairs []Pair
}

func New(pairs []Pair) (Answer, error) {
	if err := validateAnswerPairs(pairs); err != nil {
		return Answer{}, err
	}

	return Answer{
		pairs: slices.Clone(pairs),
	}, nil
}

func (a Answer) Pairs() []Pair {
	return slices.Clone(a.pairs)
}

func (a Answer) AsMap() map[PromptID]MatchID {
	pairs := make(map[PromptID]MatchID, len(a.pairs))

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
		pairs: slices.Clone(a.pairs),
	}
}
