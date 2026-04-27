package bank

import "errors"

var (
	ErrInvalid    = errors.New("ошибка банка вопросов")
	ErrDeleted    = errors.New("невозможно совершить операцию над удаленным банком вопросов")
	ErrNotDeleted = errors.New("операция восстановления возможна только над существующим банком вопросов")
)
