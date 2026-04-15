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
	ErrNilBank = errors.New("bank cannot be nil")
)

func NewSource(bankID uuid.UUID, criteria Criteria) (Source, error) {
	if bankID == uuid.Nil {
		return Source{}, ErrNilBank
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
