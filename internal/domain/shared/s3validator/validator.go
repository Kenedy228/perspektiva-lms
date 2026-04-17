package s3validator

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrEmptyS3Key   = errors.New("empty s3 key")
	ErrInvalidS3Key = errors.New("invalid s3 key")
)

func ValidateS3Key(key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrEmptyS3Key
	}

	for _, r := range key {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return ErrInvalidS3Key
		}
	}

	return nil
}
