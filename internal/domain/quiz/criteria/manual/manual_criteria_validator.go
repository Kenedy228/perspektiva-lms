package manual

import (
	"errors"
	"slices"

	"github.com/google/uuid"
)

var (
	ErrEmptyQuestions = errors.New("empty questions")
	ErrNilQuestion    = errors.New("nil question found")
	ErrInvalidCount   = errors.New("invalid criteria count")
)

func validateQuestionIDs(questionIDs []uuid.UUID) error {
	if len(questionIDs) == 0 {
		return ErrEmptyQuestions
	}

	if slices.Contains(questionIDs, uuid.Nil) {
		return ErrNilQuestion
	}

	return nil
}
