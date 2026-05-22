package source

import (
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
)

type Source struct {
	bankID   uuid.UUID
	criteria criteria.Criteria
}

func NewSource(bankID uuid.UUID, criteria criteria.Criteria) (Source, error) {
	if err := validateBankID(bankID); err != nil {
		return Source{}, err
	}

	if err := validateCriteria(criteria); err != nil {
		return Source{}, err
	}

	return Source{
		bankID:   bankID,
		criteria: criteria,
	}, nil
}

func (s Source) BankID() uuid.UUID {
	return s.bankID
}

func (s Source) Criteria() criteria.Criteria {
	return s.criteria
}

func (s Source) IsZero() bool {
	return s.bankID == uuid.Nil || s.criteria == nil
}
