package account

import (
	"fmt"

	"github.com/google/uuid"
)

func validatePersonID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: данные владельца аккаунта обязательны для заполнения", ErrInvalid)
	}

	return nil
}
