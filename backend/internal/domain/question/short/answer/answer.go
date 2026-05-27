package answer

import (
	"gitflic.ru/lms/backend/internal/domain/question"
)

type Answer struct {
	value string
}

func New(value string) Answer {
	return Answer{
		value: value,
	}
}

func (a Answer) Value() string {
	return a.value
}

func (a Answer) IsEmpty() bool {
	return len(a.value) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		value: a.value,
	}
}
