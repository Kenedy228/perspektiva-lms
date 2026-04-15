package quiz

import (
	"errors"

	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Source struct {
	id       uuid.UUID
	bankID   uuid.UUID
	criteria Criteria
}

var (
	ErrNilBank     = errors.New("bank cannot be nil")
	ErrNilCriteria = errors.New("nil criteria")
)

func NewSource(bankID uuid.UUID, criteria Criteria) (Source, error) {
	if bankID == uuid.Nil {
		return Source{}, ErrNilBank
	}

	if criteria == nil {
		return Source{}, ErrNilCriteria
	}

	id, err := utils.GenerateID()
	if err != nil {
		return Source{}, err
	}

	return Source{
		id:       id,
		bankID:   bankID,
		criteria: criteria,
	}, nil
}

func (s Source) ID() uuid.UUID {
	return s.id
}

func (s Source) BankID() uuid.UUID {
	return s.bankID
}

func (s Source) Criteria() Criteria {
	return s.criteria
}
