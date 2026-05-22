package blank

import (
	"errors"
	"regexp"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

var SinglePlaceholderRegexp = regexp.MustCompile(`^\{\{\S+?\}\}$`)

const (
	MinVariantsCount  int = 1
	MaxVariantsCount  int = 20
	VariantCharsLimit int = 1000
)

var ErrInvalid = errors.New("invalid value")

type Blank struct {
	placeholder string
	variants    []text.Text
}

func New(placeholder string, variants []text.Text) (Blank, error) {
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

func (b Blank) Variants() []text.Text {
	return slices.Clone(b.variants)
}

func (b Blank) VariantsValues() []string {
	values := make([]string, 0, len(b.variants))
	for i := range b.variants {
		values = append(values, b.variants[i].Value())
	}
	return values
}
