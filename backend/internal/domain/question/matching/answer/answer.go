package answer

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

// Pair описывает сопоставление вопроса на соответствие:
// идентификатор prompt связан с идентификатором match.
type Pair struct {
	PromptID uuid.UUID
	MatchID  uuid.UUID
}

// Answer хранит ответ пользователя для вопроса на соответствие.
type Answer struct {
	pairs []Pair
}

// New создает новый ответ и проверяет корректность пар сопоставлений.
func New(pairs []Pair) (Answer, error) {
	if err := validateAnswerPairs(pairs); err != nil {
		return Answer{}, err
	}

	return Answer{
		pairs: slices.Clone(pairs),
	}, nil
}

// Pairs возвращает копию списка сопоставлений ответа.
func (a Answer) Pairs() []Pair {
	return slices.Clone(a.pairs)
}

// AsMap возвращает сопоставления ответа в виде map[promptID]matchID.
func (a Answer) AsMap() map[uuid.UUID]uuid.UUID {
	pairs := make(map[uuid.UUID]uuid.UUID, len(a.pairs))

	for i := range a.pairs {
		pairs[a.pairs[i].PromptID] = a.pairs[i].MatchID
	}

	return pairs
}

// IsEmpty сообщает, что ответ не содержит ни одной пары сопоставлений.
func (a Answer) IsEmpty() bool {
	return len(a.pairs) == 0
}

// Clone возвращает копию ответа как question.Answer.
func (a Answer) Clone() question.Answer {
	return Answer{
		pairs: slices.Clone(a.pairs),
	}
}
