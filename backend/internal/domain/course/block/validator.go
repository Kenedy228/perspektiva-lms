package block

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/block/title"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор блока обязателен", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: название блока обязательно", ErrInvalid)
	}
	return nil
}

func validateElementIDs(ids []uuid.UUID) error {
	if len(ids) > elementsLimit {
		return fmt.Errorf("%w: количество элементов не может превышать %d", ErrInvalid, elementsLimit)
	}
	for i := range ids {
		if ids[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор элемента не может быть пустым", ErrInvalid)
		}
		if slices.Contains(ids[i+1:], ids[i]) {
			return fmt.Errorf("%w: обнаружен дублирующийся элемент", ErrInvalid)
		}
	}
	return nil
}

func validateRequiredElementID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор элемента обязателен", ErrInvalid)
	}

	return nil
}

func validateElementIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= elementsLimit {
		return fmt.Errorf("%w: превышен лимит элементов (%d)", ErrInvalid, elementsLimit)
	}

	return nil
}

func validateElementIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w: элемент уже добавлен в блок", ErrInvalid)
	}

	return nil
}
