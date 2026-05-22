package block

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/block/title"
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

func validateElementIDs(ids []uuid.UUID) error {
	if len(ids) > elementsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, elementsLimit)
	}
	for i := range ids {
		if ids[i] == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
		if slices.Contains(ids[i+1:], ids[i]) {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}
	return nil
}

func validateRequiredElementID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateElementIDsLimit(ids []uuid.UUID) error {
	if len(ids) >= elementsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, elementsLimit)
	}

	return nil
}

func validateElementIDsDuplication(target uuid.UUID, ids []uuid.UUID) error {
	if slices.Contains(ids, target) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
