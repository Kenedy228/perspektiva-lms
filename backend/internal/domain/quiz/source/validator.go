package source

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
)

func validateBankID(bankID uuid.UUID) error {
	if bankID == uuid.Nil {
		return fmt.Errorf("%w: идентификатор банка вопросов обязателен", ErrInvalid)
	}

	return nil
}

func validateCriteria(criteria criteria.Criteria) error {
	if criteria == nil {
		return fmt.Errorf("%w: критерии выборки вопросов обязательны", ErrInvalid)
	}

	return nil
}
