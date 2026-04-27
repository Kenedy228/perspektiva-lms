package content

import (
	"fmt"
)

func validateSlidersFile(key string, sizeBytes int64) error {
	if err := validateFileName(key, slidesValidExtensions); err != nil {
		return fmt.Errorf("некорректный формат файла слайдов, %w", err)
	}

	if err := validateFileSize(sizeBytes, maxSlidesFileSize); err != nil {
		return fmt.Errorf("некорректный размер файла слайдов, %w", err)
	}

	return nil
}
