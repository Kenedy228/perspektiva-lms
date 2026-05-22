package media

import (
	"errors"
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

func validateType(t Type) error {
	if !t.IsValid() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateFileForType(t Type, f file.File) error {
	if err := validateFileComplete(f); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if err := validateFileExtension(f.Extension(), t.AllowedExtensions()); err != nil {
		return fmt.Errorf("%w: invalid value (%v) (%q)", ErrInvalid, err, t.Title())
	}

	if err := validateFileSize(f.SizeBytes(), t.MaxSizeInBytes()); err != nil {
		return fmt.Errorf("%w: invalid value (%v) (%d)", ErrInvalid, err, t.MaxSizeInBytes())
	}

	return nil
}

func validateFileComplete(f file.File) error {
	if f.IsIncomplete() {
		return errors.New("invalid value")
	}

	return nil
}

func validateFileExtension(actual string, allowed []string) error {
	if !slices.Contains(allowed, actual) {
		return errors.New("invalid value")
	}

	return nil
}

func validateFileSize(actual, maxSize int64) error {
	if actual > maxSize {
		return errors.New("invalid value")
	}

	return nil
}
