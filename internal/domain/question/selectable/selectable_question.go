package selectable

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const description = "выберите один или несколько правильных вариантов ответа"
const minOptions = 2
const maxOptions = 20
const minCorrectOptions = 1

type SelectableQuestion struct {
	base.Base
	options []Option
}

func NewSelectableQuestion(params *Params) (question.Question, error) {
	base, err := base.New(&base.Params{
		Text:        params.Text,
		Description: description,
		ImageID:       params.Image,
	})

	if err != nil {
		return nil, err
	}

	if err := validateOptions(params.Options); err != nil {
		return nil, err
	}

	options := make([]Option, 0, len(params.Options))
	for text, isCorrect := range params.Options {
		option, err := NewOption(text, isCorrect)
		if err != nil {
			return nil, err
		}

		options = append(options, option)
	}

	return &SelectableQuestion{
		Base:    base,
		options: options,
	}, nil
}

func (q *SelectableQuestion) Options() []Option {
	return slices.Clone(q.options)
}

func (q *SelectableQuestion) UpdateOptions(rawOptions map[string]bool) error {
	if err := validateOptions(rawOptions); err != nil {
		return err
	}

	options := make([]Option, 0, len(rawOptions))
	for text, isCorrect := range rawOptions {
		option, err := NewOption(text, isCorrect)
		if err != nil {
			return err
		}

		options = append(options, option)
	}

	q.options = options
	q.Touch()

	return nil
}

func (q *SelectableQuestion) Type() question.Type {
	return question.TypeSelectable
}
