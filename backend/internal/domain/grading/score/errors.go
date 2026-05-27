package score

import "errors"

var (
	// ErrInvalid возвращается при недопустимом значении score.
	ErrInvalid = errors.New("некорректное значение")
)
