package slides

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/shared/file"
)

func validateFile(f file.File) error {
	if err := validateFileExtension(f); err != nil {
		return err
	}

	return nil
}

func validateFileExtension(f file.File) error {
	ext := f.Extension()

	for i := range allowedExtensions {
		if allowedExtensions[i] == ext {
			return nil
		}
	}

	return fmt.Errorf("%w, детали: файл для вложения имеет некорректный тип", ErrInvalid)
}
