package blank

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question/content"
)

type Blank struct {
	placeholder string
	variants    []content.Content
}

func New(placeholder string, variants []content.Content) (Blank, error) {
	if err := validatePlaceholder(placeholder); err != nil {
		return Blank{}, err
	}

	if err := validateVariants(variants); err != nil {
		return Blank{}, err
	}

	return Blank{
		placeholder: placeholder,
		variants:    slices.Clone(variants),
	}, nil
}

func (b Blank) Placeholder() string {
	return b.placeholder
}

func (b Blank) Variants() []content.Content {
	return slices.Clone(b.variants)
}

func (b Blank) VariantsValues() []string {
	values := make([]string, 0, len(b.variants))

	for i := range b.variants {
		values = append(values, b.variants[i].Value())
	}

	return values
}
