package typed

import "errors"

var (
	ErrInvalidText        = errors.New("некорректный текст")
	ErrInvalidBlanks      = errors.New("некорректные бланки")
	ErrInvalidPlaceholder = errors.New("некорректный заполнитель")
	ErrInvalidVariants    = errors.New("неверные варианты ответов")
)
