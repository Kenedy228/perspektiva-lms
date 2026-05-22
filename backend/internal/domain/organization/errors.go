package organization

import "errors"

// ErrInvalid is returned when organization state violates domain invariants.
var ErrInvalid = errors.New("некорректная организация")
