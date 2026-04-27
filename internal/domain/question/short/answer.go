package short

import (
	"strings"

	"gitflic.ru/lms/internal/domain/question"
)

type Answer struct {
	input string
}

func NewAnswer(params AnswerParams) question.Answer {
	return Answer{
		input: params.Input,
	}
}

func (a Answer) Input() string {
	return a.input
}

func (a Answer) IsEmpty() bool {
	return strings.TrimSpace(a.input) == ""
}

func (a Answer) Clone() question.Answer {
	return Answer{
		input: a.input,
	}
}
