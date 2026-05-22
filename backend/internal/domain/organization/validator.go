package organization

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор не существует", ErrInvalid)
	}

	return nil
}

func validateINN(inn inn.INN) error {
	if inn.IsZero() {
		return fmt.Errorf("%w: ИНН не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateName(n name.Name) error {
	if n.IsZero() {
		return fmt.Errorf("%w: наименование не может быть пустым", ErrInvalid)
	}

	return nil
}
