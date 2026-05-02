package answer

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Answer struct {
	options []AnswerOption
}

func New(options []AnswerOption) Answer {
	return Answer{
		options: slices.Clone(options),
	}
}

func (a Answer) Options() []AnswerOption {
	return slices.Clone(a.options)
}

func (a Answer) OptionsAsMap() map[uuid.UUID]struct{} {
	options := make(map[uuid.UUID]struct{}, len(a.options))

	for i := range a.options {
		options[a.options[i].OptionID] = struct{}{}
	}

	return options
}

func (a Answer) IsEmpty() bool {
	return len(a.options) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		options: slices.Clone(a.options),
	}
}
