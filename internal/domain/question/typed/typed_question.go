package typed

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const description = "заполните пропуски"
const maxPlaceholders = 20

type TypedQuestion struct {
	base.Base
	blanks []Blank
}

func New(params *Params) (question.Question, error) {
	base, err := base.New(&base.Params{
		Text:        params.Text,
		Description: description,
		Image:       params.Image,
	})

	if err != nil {
		return nil, err
	}

	if err := validatePlaceholders(params.Text, params.PlaceholdersCount, params.Blanks); err != nil {
		return nil, err
	}

	blanks := make([]Blank, 0, len(params.Blanks))
	for mark, answers := range params.Blanks {
		blank, err := NewBlank(mark, answers)
		if err != nil {
			return nil, err
		}

		blanks = append(blanks, blank)
	}

	return &TypedQuestion{
		Base:   base,
		blanks: blanks,
	}, nil
}

func (q *TypedQuestion) Blanks() []Blank {
	return slices.Clone(q.blanks)
}

func (q *TypedQuestion) ReplaceContent(text string, placeholdersCount int, rawBlanks map[string][]string) error {
	if err := validatePlaceholders(text, placeholdersCount, rawBlanks); err != nil {
		return err
	}

	blanks := make([]Blank, 0, len(rawBlanks))
	for mark, answers := range rawBlanks {
		blank, err := NewBlank(mark, answers)
		if err != nil {
			return err
		}

		blanks = append(blanks, blank)
	}

	if err := q.UpdateText(text); err != nil {
		return err
	}

	q.blanks = blanks

	return nil
}

func (q *TypedQuestion) Type() question.Type {
	return question.TypeTyped
}
