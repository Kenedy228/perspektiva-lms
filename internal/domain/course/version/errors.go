package version

import "errors"

var (
	ErrInvalid     = errors.New("ошибка версии курса")
	ErrNotEditable = errors.New("версия курса не является черновиком")
)
