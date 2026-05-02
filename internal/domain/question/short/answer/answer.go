package answer

import (
	"gitflic.ru/lms/internal/domain/question"
)

type Answer struct {
	variant AnswerVariant
}

func New(variant AnswerVariant) Answer {
	return Answer{
		variant: variant,
	}
}

func (a Answer) Variant() AnswerVariant {
	return a.variant
}

func (a Answer) VariantAsString() string {
	return a.Variant().Input
}

func (a Answer) IsEmpty() bool {
	return a.variant.Input == ""
}

func (a Answer) Clone() question.Answer {
	return Answer{
		variant: a.variant,
	}
}
