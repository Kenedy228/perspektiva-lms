package criteria

import "errors"

var (
	ErrInvalidQuestions      = errors.New("некорректная выборка вопросов")
	ErrInvalidQuestionCount = errors.New("неверное количество вопросов для случайной выборки")
)
