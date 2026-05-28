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
	if len(ids) > versionsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, versionsLimit)
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

func validateRequiredVersionID(id uuid.UUID) error {
	return validateRequiredBlockID(id)
}

func validateRequiredBlockID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateBlockIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= versionsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, versionsLimit)
	}

	return nil
}

func validateBlockIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateVersionIDs(ids []uuid.UUID) error {
	return validateBlockIDs(ids)
}

func validateVersionIDsLimit(ids []uuid.UUID) error {
	return validateBlockIDsLimit(ids)
}

func validateVersionIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	return validateBlockIDsDuplication(target, ids)
}
