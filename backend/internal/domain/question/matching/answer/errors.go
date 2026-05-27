package answer

import "errors"

var (
	// ErrInvalid возвращается при нарушении инвариантов ответа на вопрос на соответствие.
	ErrInvalid = errors.New("некорректное значение")
)
