package quiz

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNilBank = errors.New("bank cannot be nil")
)

func validateBank(bankID uuid.UUID) error {
	if bankID == uuid.Nil {
		return ErrNilBank
	}

	return nil
}
