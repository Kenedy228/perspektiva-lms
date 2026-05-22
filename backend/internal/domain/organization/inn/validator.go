package inn

import (
	"fmt"

	inn2 "github.com/Kenedy228/inn"
)

func validateValue(value string, t Type) error {
	switch t {
	case TypeIndividualEntrepreneur:
		return inn2.ValidateIndividualEntrepreneur(value)
	case TypeLegalEntity:
		return inn2.ValidateLegalEntity(value)
	case TypeNaturalPerson:
		return inn2.ValidateNaturalPerson(value)
	default:
		return fmt.Errorf("%w: неизвестный тип ИНН", ErrInvalid)
	}
}
