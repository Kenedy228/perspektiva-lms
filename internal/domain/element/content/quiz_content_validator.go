package content

import (
	"github.com/google/uuid"
)

func validateQuizContent(quizID uuid.UUID) error {
	if quizID == uuid.Nil {
		return ErrEmptyQuizID
	}

	return nil
}
