package person

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"github.com/google/uuid"
)

func validateName(n name.Name) error {
	if n.IsZero() {
		return fmt.Errorf("%w: человек должен обладать ФИО", ErrInvalid)
	}

	return nil
}

func validateProfile(prof profile.Profile) error {
	if prof.IsZero() {
		return fmt.Errorf("%w: нельзя прикрепить незаполненный профиль", ErrInvalid)
	}

	return nil
}

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор должен существовать", ErrInvalid)
	}

	return nil
}
