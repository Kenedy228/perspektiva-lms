package matching

import (
	"fmt"
)

func validatePairs(pairs []Pair) error {
	if len(pairs) < minPairs {
		return fmt.Errorf("%w, детали: пары ответов должны содержать минимум %d элементов", ErrInvalidPairs, minPairs)
	}

	if len(pairs) > maxPairs {
		return fmt.Errorf("%w, детали: пары ответов должны содержать максимум %d элементов", ErrInvalidPairs, maxPairs)
	}

	for i := range pairs {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].Prompt() == pairs[j].Prompt() {
				return fmt.Errorf("%w, детали: пары ответов не могут содержать дубликаты в столбце понятий", ErrInvalidPairs)
			}

			if pairs[i].Content() == pairs[j].Content() {
				return fmt.Errorf("%w, детали: пары ответов не могут содержать дубликаты в столбце определений", ErrInvalidPairs)
			}
		}
	}

	return nil
}
