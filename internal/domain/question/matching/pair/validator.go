package pair

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/question/content"
)

func validatePrompt(prompt content.Content) error {
	if err := validatePromptContentType(prompt); err != nil {
		return err
	}

	if err := validatePromptLimit(prompt); err != nil {
		return err
	}

	return nil
}

func validateMatch(match content.Content) error {
	if !match.IsText() {
		return nil
	}

	if err := validateMatchAsTextLimit(match); err != nil {
		return err
	}

	return nil
}

func validatePromptContentType(prompt content.Content) error {
	if !prompt.IsText() {
		return fmt.Errorf("%w, детали: определение должно быть текстового вида", ErrInvalid)
	}

	return nil
}

func validatePromptLimit(prompt content.Content) error {
	if utf8.RuneCountInString(prompt.Value()) > promptCharsLimit {
		return fmt.Errorf("%w, детали: определение должно содержать не более %d символов", ErrInvalid, promptCharsLimit)
	}

	return nil
}

func validateMatchAsTextLimit(match content.Content) error {
	if utf8.RuneCountInString(match.Value()) > matchAsTextCharsLimit {
		return fmt.Errorf("%w, детали: текстовое соответствие должно содержать не более %d символов", ErrInvalid, matchAsTextCharsLimit)
	}

	return nil
}
