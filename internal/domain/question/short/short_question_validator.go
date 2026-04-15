package short

import (
	"errors"
	"strings"
)

var (
	ErrNoAnswers       = errors.New("no answers provided")
	ErrEmptyAnswer     = errors.New("answer is empty")
	ErrTooManyAnswers  = errors.New("too many answers")
	ErrDuplicateAnswer = errors.New("duplicate answer found")
)

func validateAnswers(answers []string, allowDuplicates bool) error {
	if len(answers) == 0 {
		return ErrNoAnswers
	}

	if len(answers) > maxAnswers {
		return ErrTooManyAnswers
	}

	visited := make(map[string]struct{}, len(answers))

	for i := range answers {
		if strings.TrimSpace(answers[i]) == "" {
			return ErrEmptyAnswer
		}

		if _, ok := visited[answers[i]]; ok && !allowDuplicates {
			return ErrDuplicateAnswer
		}

		visited[answers[i]] = struct{}{}
	}

	return nil
}
