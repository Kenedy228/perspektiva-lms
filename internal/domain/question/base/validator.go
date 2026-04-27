package base

import (
	"fmt"
	"strings"
)

func validateText(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("%w, детали: текст вопроса должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}
