package file

import (
	"fmt"
	"path"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
)

func validateFile(key string, sizeBytes int64) error {
	if err := validateFileName(key); err != nil {
		return err
	}

	if err := validateFileSize(sizeBytes); err != nil {
		return err
	}

	return nil
}

func validateFileName(key string) error {
	if err := s3validator.ValidateS3Key(key); err != nil {
		return err
	}

	filename := path.Base(key)
	if filename == "" || filename == "." {
		return fmt.Errorf("%w, детали: некорректный путь к файлу", ErrInvalid)
	}

	ext := strings.ToLower(path.Ext(filename))

	nameWithoutExt := strings.TrimSuffix(filename, ext)
	if nameWithoutExt == "" {
		return fmt.Errorf("%w, детали: некорректный путь к файлу", ErrInvalid)
	}

	return nil
}

func validateFileSize(sizeBytes int64) error {
	if sizeBytes <= 0 {
		return fmt.Errorf("%w, детали: размер файла не может быть отрицательным", ErrInvalid)
	}

	if sizeBytes > maxSizeBytes {
		return fmt.Errorf("%w, детали: размер файла должен быть не более %d байт", ErrInvalid, maxSizeBytes)
	}

	return nil
}
