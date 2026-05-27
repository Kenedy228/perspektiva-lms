package attempt

import "errors"

var (
	ErrInvalid       = errors.New("некорректное значение")
	ErrStateConflict = errors.New("конфликт состояния")
	ErrNotFound      = errors.New("не найдено")
)
