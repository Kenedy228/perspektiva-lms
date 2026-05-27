package matching

import "errors"

var (
	// ErrInvalid сигнализирует о нарушении доменных инвариантов вопроса matching.
	ErrInvalid = errors.New("некорректное значение")
)
