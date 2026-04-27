package matching

import "errors"

var (
	ErrInvalidPairs  = errors.New("неправильные пары ответов")
	ErrInvalidPrompt = errors.New("неверное определение")
)
