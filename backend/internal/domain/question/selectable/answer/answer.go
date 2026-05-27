package answer

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

type Answer struct {
	optionIDs []uuid.UUID
}

func New(optionIDs []uuid.UUID) (Answer, error) {
	if err := validateOptionIDs(optionIDs); err != nil {
		return Answer{}, err
	}

	return Answer{
		optionIDs: slices.Clone(optionIDs),
	}, nil
}

func (a Answer) OptionIDs() []uuid.UUID {
	return slices.Clone(a.optionIDs)
}

func (a Answer) OptionIDSet() map[uuid.UUID]struct{} {
	options := make(map[uuid.UUID]struct{}, len(a.optionIDs))

	for i := range a.optionIDs {
		options[a.optionIDs[i]] = struct{}{}
	}

	return options
}

func (a Answer) IsEmpty() bool {
	return len(a.optionIDs) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		optionIDs: slices.Clone(a.optionIDs),
	}
}
