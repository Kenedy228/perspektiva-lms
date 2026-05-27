package matching

import "errors"

var (
	// ErrInvalidQuestionType возвращается, когда checker получил вопрос неподдерживаемого типа.
	ErrInvalidQuestionType = errors.New("неверный тип вопроса")
	// ErrInvalidAnswerType возвращается, когда checker получил ответ неподдерживаемого типа.
	ErrInvalidAnswerType = errors.New("неверный тип ответа")
)
