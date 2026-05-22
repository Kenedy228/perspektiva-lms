package pair

import (
	"fmt"
)

func validatePrompt(prompt Prompt) error {
	if prompt.IsIncomplete() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateMatch(match Match) error {
	if match.IsIncomplete() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
