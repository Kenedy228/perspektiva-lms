package typed

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question/option"
)

type Blank struct {
	placeholder string
	answers     []option.ContentOption
}

func NewBlank(params BlankParams) (Blank, error) {
	if err := validatePlaceholder(params.Placeholder); err != nil {
		return Blank{}, err
	}

	if err := validateAnswers(params.Answers); err != nil {
		return Blank{}, err
	}

	cAnswers := slices.Clone(params.Answers)

	return Blank{
		placeholder: params.Placeholder,
		answers:     cAnswers,
	}, nil
}

func (b Blank) Placeholder() string {
	return b.placeholder
}

func (b Blank) Answers() []option.ContentOption {
	return slices.Clone(b.answers)
}
