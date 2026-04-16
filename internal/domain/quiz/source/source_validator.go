package source

import (
	"errors"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
)

var (
	ErrNilBank     = errors.New("bank cannot be nil")
	ErrNilCriteria = errors.New("nil criteria")
)

func validateBankID(bankID uuid.UUID) error {
	if bankID == uuid.Nil {
		return ErrNilBank
	}

	return nil
}

func validateCriteria(criteria criteria.Criteria) error {
	if criteria == nil {
		return ErrNilCriteria
	}

	return nil
}
