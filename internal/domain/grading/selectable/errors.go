package selectable

import "errors"

var (
	ErrInvalidQuestionType = errors.New("неподдерживаемый формат вопроса")
	ErrInvalidAnswerType   = errors.New("неподдерживаемый формат ответа")
)
