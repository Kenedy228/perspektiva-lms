package selectable

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Answer struct {
	optionIDs []uuid.UUID
}

func NewAnswer(params AnswerParams) question.Answer {
	return Answer{
		optionIDs: params.OptionIDs,
	}
}

func (a Answer) OptionIDs() []uuid.UUID {
	return slices.Clone(a.optionIDs)
}

func (a Answer) IsEmpty() bool {
	return len(a.optionIDs) == 0
}

func (a Answer) Clone() question.Answer {
	cOptionIDs := slices.Clone(a.optionIDs)
	return Answer{
		optionIDs: cOptionIDs,
	}
}
