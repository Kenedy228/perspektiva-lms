package person

import "errors"

var (
	ErrEmptyFirstName = errors.New("firstName cannot be empty")
	ErrEmptyLastName  = errors.New("lastName cannot be empty")
)
