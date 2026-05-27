package grading

import "errors"

var (
	ErrInvalidQuestionType = errors.New("неверный тип вопроса")
	ErrInvalidAnswerType   = errors.New("неверный тип ответа")
	ErrNilQuestion         = errors.New("вопрос не может быть nil")
	ErrNilAnswer           = errors.New("ответ не может быть nil")
)
