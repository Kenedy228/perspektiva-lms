package selectable

import (
	"errors"
	"strings"
)

var (
	ErrEmptyOptionText         = errors.New("empty option text")
	ErrEmptyOptions            = errors.New("empty options")
	ErrNotEnoughOptions        = errors.New("not enough options")
	ErrNoCorrectOption         = errors.New("no correct option provided")
	ErrNotEnoughCorrectOptions = errors.New("not enough correct options")
	ErrTooManyOptions          = errors.New("too many options")
)

func validateOptionText(text string) error {
	if strings.TrimSpace(text) == "" {
		return ErrEmptyOptionText
	}

	return nil
}

func validateOptions(options map[string]bool) error {
	if len(options) == 0 {
		return ErrEmptyOptions
	}

	if len(options) < minOptions {
		return ErrNotEnoughOptions
	}

	if len(options) > maxOptions {
		return ErrTooManyOptions
	}

	correctCount := 0
	for _, isCorrect := range options {
		if isCorrect {
			correctCount++
		}
	}

	if correctCount == 0 {
		return ErrNoCorrectOption
	}

	if correctCount < minCorrectOptions {
		return ErrNotEnoughCorrectOptions
	}

	return nil
}
