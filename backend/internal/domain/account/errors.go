package account

import "errors"

var (
	ErrInvalid       = errors.New("некорректный аккаунт")
	ErrAlreadyActive = errors.New("аккаунт уже активен")
)
