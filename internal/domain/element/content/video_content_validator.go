package content

import "fmt"

func validateVideoFile(key string, sizeBytes int64) error {
	if err := validateFileName(key, videoValidExtensions); err != nil {
		return fmt.Errorf("некорректный формат видео, %w", err)
	}

	if err := validateFileSize(sizeBytes, maxVideoFileSize); err != nil {
		return fmt.Errorf("некорректный размер видео, %w", err)
	}

	return nil
}
