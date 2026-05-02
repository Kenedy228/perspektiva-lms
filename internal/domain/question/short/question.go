package short

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/short/variant"
	"gitflic.ru/lms/internal/domain/question/title"
)

type Question struct {
	*base.Base
	variants []variant.Variant
}

func New(t title.Title, variants []variant.Variant) (*Question, error) {
	base, err := base.New(t)
	if err != nil {
		return nil, err
	}

	if err := validateVariants(variants); err != nil {
		return nil, err
	}

	return &Question{
		Base:     base,
		variants: slices.Clone(variants),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Variants() []variant.Variant {
	return slices.Clone(q.variants)
}

func (q *Question) Type() question.Type {
	return question.TypeShort
}

func (q *Question) ChangeVariants(variants []variant.Variant) error {
	if err := validateVariants(variants); err != nil {
		return err
	}

	q.variants = slices.Clone(variants)
	return nil
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:     q.Base.Clone(),
		variants: slices.Clone(q.variants),
	}
}
