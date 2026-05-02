package course

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

func validateRequiredVersionID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: передан несуществующий идентификатор версии", ErrInvalid)
	}

	return nil
}

func validateVersionIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= versionsLimit {
		return fmt.Errorf("%w, детали: количество версий в курсе не должно превышать %d штук", ErrInvalid, versionsLimit)
	}

	return nil
}

func validateVersionIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w, детали: курс не должен содержать дубликаты версий", ErrInvalid)
	}

	return nil
}
