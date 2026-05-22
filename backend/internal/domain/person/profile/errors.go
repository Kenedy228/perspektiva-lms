package profile

import "errors"

// ErrInvalid is returned when profile state violates domain invariants.
var ErrInvalid = errors.New("некорректный профиль")
