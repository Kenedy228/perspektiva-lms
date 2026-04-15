package quiz

import (
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Strategy interface {
}

type Source struct {
	id       uuid.UUID
	bank     uuid.UUID
	strategy Strategy
}

func New(bankID uuid.UUID, strategy Strategy) (Source, error) {
	if err := validateBank(bankID); err != nil {
		return Source{}, err
	}

	id, err := utils.GenerateID()

	if err != nil {
		return Source{}, err
	}

	return Source{
		id:       id,
		bank:     bankID,
		strategy: strategy,
	}, nil
}


