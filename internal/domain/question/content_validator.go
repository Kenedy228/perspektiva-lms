package question

import (
	"errors"
	"fmt"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/s3validator"
)

var (
	ErrInvalidContent = errors.New("некорректный контент вопроса")
)

func validateContentType(cType ContentType) error {
	if !cType.IsValid() {
		return fmt.Errorf("%w, детали: некорректный вид контента", ErrInvalidContent)
	}

	return nil
}

func validateValue(cType ContentType, val string) error {
	switch cType {
	case ContentTypeText:
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("%w, детали: для текстового типа контента содержимое должно иметь хотя бы один непробельный символ", ErrInvalidContent)
		}
	case ContentTypeImage, ContentTypeAudio:
		if err := s3validator.ValidateS3Key(val); err != nil {
			return fmt.Errorf("%w, детали: для контента вида изображение или аудио содержимое должно быть в виде корректной ссылки на облачное хранилище", ErrInvalidContent)
		}
	}

	return nil
}
