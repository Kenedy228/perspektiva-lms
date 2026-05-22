package profile

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
)

func validateSNILS(s snils.SNILS) error {
	if s.IsZero() {
		return fmt.Errorf("%w: СНИЛС должен быть указан для профиля", ErrInvalid)
	}

	return nil
}

func validateDateOfBirth(dob dob.DateOfBirth) error {
	if dob.IsZero() {
		return fmt.Errorf("%w: дата рождения должна быть указана для профиля", ErrInvalid)
	}

	return nil
}
