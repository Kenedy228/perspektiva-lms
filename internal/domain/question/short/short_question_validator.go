package short

import (
	"errors"

	"gitflic.ru/lms/internal/domain/question/option"
)

var (
	ErrNoAnswers           = errors.New("no answers provided")
	ErrTooManyAnswers      = errors.New("too many answers")
	ErrDuplicateAnswer     = errors.New("duplicate answer found")
	ErrInvalidAnswerFormat = errors.New("invalid answer format")
)

func validateAnswers(answers []option.ContentOption) error {
	if len(answers) == 0 {
		return ErrNoAnswers
	}

	if len(answers) > maxAnswers {
		return ErrTooManyAnswers
	}

	visited := make(map[option.ContentOption]struct{}, len(answers))

	for i := range answers {
		if answers[i].ContentType() != option.ContentTypeText {
			return ErrInvalidAnswerFormat
		}

		if _, ok := visited[answers[i]]; ok {
			return ErrDuplicateAnswer
		}

		visited[answers[i]] = struct{}{}
	}

	return nil
}
