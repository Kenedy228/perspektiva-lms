package quiz

import (
	"fmt"

	"github.com/google/uuid"
)

func validateQuizID(quizID uuid.UUID) error {
	if err := validateRequiredQuizID(quizID); err != nil {
		return err
	}

	return nil
}

func validateRequiredQuizID(quizID uuid.UUID) error {
	if quizID == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
