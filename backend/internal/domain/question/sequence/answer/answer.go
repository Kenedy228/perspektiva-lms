package answer

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
)

type Answer struct {
	optionIDs []OptionID
}

func New(optionIDs []OptionID) (Answer, error) {
	if err := validateOptionIDs(optionIDs); err != nil {
		return Answer{}, err
	}

	return Answer{
		optionIDs: slices.Clone(optionIDs),
	}, nil
}

func (a Answer) OptionIDs() []OptionID {
	return slices.Clone(a.optionIDs)
}

func (a Answer) IsEmpty() bool {
	return len(a.optionIDs) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		optionIDs: slices.Clone(a.optionIDs),
	}
}
