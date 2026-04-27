package quiz

import "errors"

var (
	ErrInvalidSource          = errors.New("неверный источник вопросов")
	ErrInvalid                = errors.New("ошибка квиза")
	ErrSourceSizeExceeded     = errors.New("достигнут лимит количества источников")
	ErrDuplicateBankID        = errors.New("дубликат банка вопросов в квизе")
	ErrCannotRemoveLastSource = errors.New("ошибка удаления источника")
)
