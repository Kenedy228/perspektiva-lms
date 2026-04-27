package sequence

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Answer struct {
	elementIDs []uuid.UUID
}

func NewAnswer(params AnswerParams) question.Answer {
	return Answer{
		elementIDs: params.ElementIDs,
	}
}

func (a Answer) IDs() []uuid.UUID {
	return slices.Clone(a.elementIDs)
}

func (a Answer) IsEmpty() bool {
	return len(a.elementIDs) == 0
}

func (a Answer) Clone() question.Answer {
	cIds := slices.Clone(a.elementIDs)
	return Answer{
		elementIDs: cIds,
	}
}
