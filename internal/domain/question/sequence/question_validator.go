package sequence

import (
	"fmt"
)

func validateElements(elements []Element) error {
	if len(elements) < minElements {
		return fmt.Errorf("%w, детали: количество элементов должно быть не менее %d штук", ErrInvalidElements, minElements)
	}

	if len(elements) > maxElements {
		return fmt.Errorf("%w, детали: количество элементов должно быть не более %d штук", ErrInvalidElements, maxElements)
	}

	for i := range elements {
		for j := i + 1; j < len(elements); j++ {
			if elements[i].Equal(elements[j]) {
				return fmt.Errorf("%w, детали: дубликаты элементов не допустимы", ErrInvalidElements)
			}
		}
	}

	return nil
}
