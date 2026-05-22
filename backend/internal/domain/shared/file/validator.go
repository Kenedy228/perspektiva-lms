package file

import (
	"errors"
	"fmt"
	"io/fs"
)

func validateFile(fileName string, sizeBytes int64) error {
	if err := validateFileName(fileName); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if err := validateFileSize(sizeBytes); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	return nil
}

func validateFileName(fileName string) error {
	if !fs.ValidPath(fileName) {
		return errors.New("invalid value")
	}

	return nil
}

func validateFileSize(sizeBytes int64) error {
	if sizeBytes <= 0 {
		return errors.New("invalid value")
	}

	return nil
}
