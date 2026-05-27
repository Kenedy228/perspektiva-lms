package short

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
)

type Question struct {
	*base.Base
	variants []variant.Variant
}

func New(b *base.Base, variants []variant.Variant) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validateVariants(variants); err != nil {
		return nil, err
	}

	return &Question{
		Base:     b,
		variants: slices.Clone(variants),
	}, nil
}

func Restore(b *base.Base, variants []variant.Variant) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validateVariants(variants); err != nil {
		return nil, err
	}

	return &Question{
		Base:     b,
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
