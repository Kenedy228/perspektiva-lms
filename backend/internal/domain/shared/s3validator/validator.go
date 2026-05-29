// Пакет s3validator предоставляет функционал по валидации s3 ключей.
package s3validator

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrEmptyS3Key      = errors.New("ключ объекта хранилища не может быть пустым")
	ErrInvalidS3Key    = errors.New("ключ объекта хранилища содержит недопустимые символы")
	ErrTooLongS3Key    = errors.New("ключ объекта хранилища превышает допустимую длину")
	ErrUnsafeS3KeyPath = errors.New("ключ объекта хранилища содержит небезопасный путь")
)

// Максимальный размер ключа в байтах.
const maxS3KeyBytes = 1024

// Функция валидации s3 ключа, реализована в соответствии с документаций amazon s3.
func ValidateS3Key(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrEmptyS3Key
	}

	if len(key) > maxS3KeyBytes {
		return ErrTooLongS3Key
	}

	if !utf8.ValidString(key) {
		return ErrInvalidS3Key
	}

	if strings.Contains(key, "//") || strings.Contains(key, "\\") {
		return ErrInvalidS3Key
	}

	if strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") {
		return ErrInvalidS3Key
	}

	parts := strings.Split(key, "/")
	for i := range parts {
		if parts[i] == "" || parts[i] == "." || parts[i] == ".." {
			return ErrUnsafeS3KeyPath
		}

		for _, r := range parts[i] {
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				return ErrInvalidS3Key
			}
		}
	}

	return nil
}
