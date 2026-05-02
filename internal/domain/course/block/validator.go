package block

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

func validateRequiredElementID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: передан несуществующий идентификатор элемента", ErrInvalid)
	}

	return nil
}

func validateElementIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= elementsLimit {
		return fmt.Errorf("%w, детали: количество элементов в блоке не должно превышать %d штук", ErrInvalid, elementsLimit)
	}

	return nil
}

func validateElementIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w, детали: блок не должен содержать дубликаты элементов (разрешаются копии)", ErrInvalid)
	}

	return nil
}
