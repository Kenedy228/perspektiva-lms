package selectable

import (
	"fmt"
)

func validateOptions(options []Option) error {
	if len(options) < minOptions {
		return fmt.Errorf("%w, детали: количество вариантов ответов не должно быть меньше %d штук", ErrInvalidOptions, minOptions)
	}

	if len(options) > maxOptions {
		return fmt.Errorf("%w, детали: количество вариантов ответов не должно быть больше %d штук", ErrInvalidOptions, maxOptions)
	}

	for i := range options {
		for j := i + 1; j < len(options); j++ {
			if options[i].Equal(options[j]) {
				return fmt.Errorf("%w, детали: варианты ответов не должны содержать одинаковый контент", ErrInvalidOptions)
			}
		}
	}

	correctCount := countCorrect(options)

	if correctCount < minCorrect {
		return fmt.Errorf("%w, детали: варианты ответов должны содержать минимум %d правильных ответов", ErrInvalidOptions, minCorrect)
	}

	return nil
}
