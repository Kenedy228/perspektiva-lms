package attempt

import "errors"

var (
	ErrInvalid       = errors.New("некорректная попытка")
	ErrStateConflict = errors.New("конфликт состояния попытки")
	ErrNotFound      = errors.New("вопрос не найден в данной попытке")
)
