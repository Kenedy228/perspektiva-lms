package course

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/title"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор курса обязателен", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: название курса обязательно", ErrInvalid)
	}
	return nil
}

func validateBlockIDs(ids []uuid.UUID) error {
	if len(ids) > blocksLimit {
		return fmt.Errorf("%w: количество блоков не может превышать %d", ErrInvalid, blocksLimit)
	}
	for i := range ids {
		if ids[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор блока не может быть пустым", ErrInvalid)
		}
	}
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for i := range ids {
		if _, ok := seen[ids[i]]; ok {
			return fmt.Errorf("%w: обнаружен дублирующийся блок", ErrInvalid)
		}
		seen[ids[i]] = struct{}{}
	}
	return nil
}

func validateRequiredBlockID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор блока обязателен", ErrInvalid)
	}

	return nil
}

func validateBlockIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= blocksLimit {
		return fmt.Errorf("%w: превышен лимит блоков (%d)", ErrInvalid, blocksLimit)
	}

	return nil
}

func validateBlockIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w: блок уже добавлен в курс", ErrInvalid)
	}

	return nil
}
