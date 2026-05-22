package person

import "errors"

// ErrInvalid is returned when person state violates domain invariants.
var ErrInvalid = errors.New("некорректный человек")
