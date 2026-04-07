package account

import "errors"

var (
	ErrEmptyLogin        = errors.New("login can't be empty")
	ErrEmptyPasswordHash = errors.New("password hash can't be empty")
)
