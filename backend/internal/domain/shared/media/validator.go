package media

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/shared/file"
)

func validateType(t Type) error {
	if !t.IsValid() {
		return fmt.Errorf("%w: неизвестный тип медиафайла", ErrInvalid)
	}

	return nil
}

func validateFileForType(t Type, f file.File) error {
	if err := validateFileComplete(f); err != nil {
		return err
	}

	if err := validateFileExtension(f.Extension(), t.AllowedExtensions()); err != nil {
		return fmt.Errorf("%w: %v для типа %q", ErrInvalid, err, t.Title())
	}

	if err := validateFileSize(f.SizeBytes(), t.MaxSizeInBytes()); err != nil {
		return fmt.Errorf("%w: %v (лимит %d байт)", ErrInvalid, err, t.MaxSizeInBytes())
	}

	return nil
}

func validateFileComplete(f file.File) error {
	if f.IsIncomplete() {
		return fmt.Errorf("%w: медиафайл не заполнен", ErrInvalid)
	}

	return nil
}

func validateFileExtension(actual string, allowed []string) error {
	if !slices.Contains(allowed, actual) {
		return fmt.Errorf("недопустимое расширение файла %q", actual)
	}

	return nil
}

func validateFileSize(actual, maxSize int64) error {
	if actual > maxSize {
		return fmt.Errorf("превышен максимальный размер файла")
	}

	return nil
}
