package typed

import (
	"maps"

	"gitflic.ru/lms/internal/domain/question"
)

type Answer struct {
	inputs map[string]string
}

func NewAnswer(params AnswerParams) question.Answer {
	return Answer{
		inputs: params.Inputs,
	}
}

func (a Answer) Inputs() map[string]string {
	return maps.Clone(a.inputs)
}

func (a Answer) IsEmpty() bool {
	return len(a.inputs) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		inputs: a.Inputs(),
	}
}
