package matching

import (
	"fmt"
	"strings"
)

func validatePrompt(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return fmt.Errorf("%w, детали: определение должно содержать как минимум один непробельный символ", ErrInvalidPrompt)
	}

	return nil
}
