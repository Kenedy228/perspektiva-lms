package course

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/title"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateBlockIDs(ids []uuid.UUID) error {
	if len(ids) > blocksLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, blocksLimit)
	}
	for i := range ids {
		if ids[i] == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}
	seen := make(map[uuid.UUID]struct{}, len(ids))
	for i := range ids {
		if _, ok := seen[ids[i]]; ok {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
		seen[ids[i]] = struct{}{}
	}
	return nil
}

func validateRequiredBlockID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateBlockIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= blocksLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, blocksLimit)
	}

	return nil
}

func validateBlockIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
