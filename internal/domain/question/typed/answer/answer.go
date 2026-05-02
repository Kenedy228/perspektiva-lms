package answer

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
)

type Answer struct {
	blanks []AnswerBlank
}

func New(blanks []AnswerBlank) Answer {
	return Answer{
		blanks: blanks,
	}
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
