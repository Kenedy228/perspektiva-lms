package answer

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/question"
)

const InputCharsLimit int = 1000

var ErrInvalid = errors.New("invalid value")

type Answer struct {
	input string
}

func New(input string) (Answer, error) {
	if utf8.RuneCountInString(input) > InputCharsLimit {
		return Answer{}, fmt.Errorf("%w: invalid value (%d)", ErrInvalid, InputCharsLimit)
	}

	return Answer{
		input: input,
	}, nil
}

func (a Answer) Value() string {
	return a.input
}

func (a Answer) IsEmpty() bool {
	return len(a.input) == 0
}

func (a Answer) Clone() question.Answer {
	return Answer{
		input: a.input,
	}
}
