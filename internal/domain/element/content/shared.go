package content

import (
	"path"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
)

func validateFileName(key string, extensions []string) error {
	if err := s3validator.ValidateS3Key(key); err != nil {
		return err
	}

	filename := path.Base(key)

	if filename == "" || filename == "." {
		return ErrInvalidFormat
	}

	ext := strings.ToLower(path.Ext(filename))
	if !isValidFileExtension(ext, extensions) {
		return ErrInvalidFormat
	}

	nameWithoutExt := strings.TrimSuffix(filename, ext)
	if nameWithoutExt == "" {
		return ErrInvalidFormat
	}

	return nil
}

func validateFileSize(sizeBytes, max int64) error {
	if sizeBytes <= 0 {
		return ErrEmptyFileSize
	}

	if sizeBytes > max {
		return ErrTooLargeFile
	}

	return nil
}

func isValidFileExtension(target string, extensions []string) bool {
	for i := range extensions {
		if extensions[i] == target {
			return true
		}
	}

	return false
}
