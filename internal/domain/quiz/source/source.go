package source

import (
	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Source struct {
	id       uuid.UUID
	bankID   uuid.UUID
	criteria criteria.Criteria
}

func NewSource(params Params) (Source, error) {
	if err := validateBankID(params.BankID); err != nil {
		return Source{}, err
	}

	if err := validateCriteria(params.Criteria); err != nil {
		return Source{}, err
	}

	id, err := utils.GenerateID()
	if err != nil {
		return Source{}, err
	}

	return Source{
		id:       id,
		bankID:   params.BankID,
		criteria: params.Criteria,
	}, nil
}

func (s Source) ID() uuid.UUID {
	return s.id
}

func (s Source) BankID() uuid.UUID {
	return s.bankID
}

func (s Source) Criteria() criteria.Criteria {
	return s.criteria
}

func (s Source) Equal(other Source) bool {
	return s.id == other.ID()
}
