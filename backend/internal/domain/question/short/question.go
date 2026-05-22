package short

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type Question struct {
	*base.Base
	variants []variant.Variant
}

func New(t title.Title, variants []variant.Variant) (*Question, error) {
	if err := validateVariants(variants); err != nil {
		return nil, err
	}

	b, err := base.New(t)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:     b,
		variants: slices.Clone(variants),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, att *attachment.Attachment, variants []variant.Variant) (*Question, error) {
	if err := validateVariants(variants); err != nil {
		return nil, err
	}

	b, err := base.Restore(id, t, att)
	if err != nil {
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
