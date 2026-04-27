package short

import (
	"slices"
	"strings"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

type Question struct {
	*base.Base
	variants []question.Content
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateVariants(params.Variants); err != nil {
		return nil, err
	}

	cVariants := slices.Clone(params.Variants)

	return &Question{
		Base:     base,
		variants: cVariants,
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Variants() []question.Content {
	return slices.Clone(q.variants)
}

func (q *Question) UpdateVariants(variants []question.Content) error {
	if err := validateVariants(variants); err != nil {
		return err
	}

	cVariants := slices.Clone(variants)
	q.variants = cVariants
	q.Touch()
	return nil
}

func (q *Question) Type() question.Type {
	return question.TypeShort
}

func (q *Question) CheckAnswer(answer question.Answer) bool {
	cast, ok := answer.(Answer)
	if !ok {
		return false
	}

	if cast.IsEmpty() {
		return false
	}

	for i := range q.variants {
		if strings.EqualFold(q.variants[i].Value(), cast.Input()) {
			return true
		}
	}

	return false
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:     q.Base.Clone(),
		variants: q.Variants(),
	}
}
