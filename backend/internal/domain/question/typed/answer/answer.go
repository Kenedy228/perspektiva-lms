package answer

import (
	"errors"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
)

var ErrInvalid = errors.New("invalid value")

type Answer struct {
	blanks []AnswerBlank
}

func New(blanks []AnswerBlank) (Answer, error) {
	if err := validateBlanks(blanks); err != nil {
		return Answer{}, err
	}

	return Answer{
		blanks: slices.Clone(blanks),
	}, nil
}

func (a Answer) Blanks() []AnswerBlank {
	return slices.Clone(a.blanks)
}

func (a Answer) BlanksAsMap() map[string]string {
	blanks := make(map[string]string, len(a.blanks))

	for i := range a.blanks {
		blanks[a.blanks[i].Placeholder] = a.blanks[i].Variant
	}

	return blanks
}

func (a Answer) IsEmpty() bool {
	return len(a.blanks) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		blanks: slices.Clone(a.blanks),
	}
}
