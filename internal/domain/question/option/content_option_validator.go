package option

import (
	"errors"
	"fmt"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
)

var (
	ErrInvalidContentType = errors.New("invalid content type")
	ErrInvalidValue       = errors.New("invalid value")
)

func validateContentType(cType ContentType) error {
	if !cType.IsValid() {
		return ErrInvalidContentType
	}

	return nil
}

func validateValue(cType ContentType, val string) error {
	switch cType {
	case ContentTypeText:
		if strings.TrimSpace(val) == "" {
			return ErrInvalidValue
		}
	case ContentTypeImage, ContentTypeAudio:
		if err := s3validator.ValidateS3Key(val); err != nil {
			return fmt.Errorf("%w, details: %w", ErrInvalidValue, err)
		}
	}

	return nil
}
