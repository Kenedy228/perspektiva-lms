package quiz

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
)

func validateBankID(bankID uuid.UUID) error {
	if bankID == uuid.Nil {
		return fmt.Errorf("%w, детали: указан несуществующий банк вопросов", ErrInvalidSource)
	}

	return nil
}

func validateCriteria(criteria criteria.Criteria) error {
	if criteria == nil {
		return fmt.Errorf("%w, детали: не выбран критерий выборки", ErrInvalidSource)
	}

	return nil
}
