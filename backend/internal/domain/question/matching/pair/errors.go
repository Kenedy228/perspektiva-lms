package pair

import "errors"

var (
	// ErrInvalid возвращается при нарушении инвариантов pair value-объектов.
	ErrInvalid = errors.New("некорректное значение")
)
