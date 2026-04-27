package typed

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
)

type Blank struct {
	placeholder string
	variants    []question.Content
}

func NewBlank(params BlankParams) (Blank, error) {
	if err := validatePlaceholder(params.Placeholder); err != nil {
		return Blank{}, err
	}

	if err := validateVariants(params.Variants); err != nil {
		return Blank{}, err
	}

	cAnswers := slices.Clone(params.Variants)

	return Blank{
		placeholder: params.Placeholder,
		variants:    cAnswers,
	}, nil
}

func (b Blank) Placeholder() string {
	return b.placeholder
}

func (b Blank) Variants() []question.Content {
	return slices.Clone(b.variants)
}
