package source

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
)

func validateBankID(bankID uuid.UUID) error {
	if bankID == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateCriteria(criteria criteria.Criteria) error {
	if criteria == nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
