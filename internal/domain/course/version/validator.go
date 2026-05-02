package version

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

func validateRequiredBlockID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: передан несуществующий идентификатор блока", ErrInvalid)
	}

	return nil
}

func validateBlockIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= blocksLimit {
		return fmt.Errorf("%w, детали: количество блоков в версии не должно превышать %d штук", ErrInvalid, blocksLimit)
	}

	return nil
}

func validateBlockIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w, детали: версия не должна содержать дубликаты блоков", ErrInvalid)
	}

	return nil
}
