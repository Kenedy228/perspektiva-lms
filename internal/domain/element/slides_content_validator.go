package element

import (
	"fmt"
	"path"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
)

const maxSlidesFileSizeBytes int64 = 100 * 1024 * 1024 // 100 MB

var (
	validExtensions = []string{".pptx"}
)

func validateSlidersFile(key string, sizeBytes int64) error {
	if err := validateSlidersFileName(key); err != nil {
		return err
	}

	if err := validateSlidersFileSize(sizeBytes); err != nil {
		return err
	}

	return nil
}

func validateSlidersFileName(key string) error {
	if err := s3validator.ValidateS3Key(key); err != nil {
		return err
	}

	filename := path.Base(key)

	if filename == "" || filename == "." {
		return fmt.Errorf("%w, детали: неккоректный формат файла слайдов", ErrInvalid)
	}

	ext := strings.ToLower(path.Ext(filename))
	if !isValidFileExtension(ext) {
		return fmt.Errorf("%w, детали: неккоректный формат файла слайдов", ErrInvalid)

	}

	nameWithoutExt := strings.TrimSuffix(filename, ext)
	if nameWithoutExt == "" {
		return fmt.Errorf("%w, детали: неккоректный формат файла слайдов", ErrInvalid)
	}

	return nil
}

func isValidFileExtension(ext string) bool {
	for i := range validExtensions {
		if validExtensions[i] == ext {
			return true
		}
	}

	return false
}

func validateSlidersFileSize(sizeBytes int64) error {
	if sizeBytes <= 0 {
		return fmt.Errorf("%w, детали: пустой размер файла слайдов", ErrInvalid)
	}

	if sizeBytes > maxSlidesFileSizeBytes {
		return fmt.Errorf("%w, детали: размер файла слайдов превышает допустимый размер", ErrInvalid)
	}

	return nil
}
