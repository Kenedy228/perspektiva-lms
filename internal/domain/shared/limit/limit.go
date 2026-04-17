package limit

import "errors"

type Limit int

var (
	ErrInvalidLimit = errors.New("invalid limit")
)

func (l Limit) Validate() error {
	if l < 0 {
		return ErrInvalidLimit
	}

	return nil
}

func (l Limit) IsInfinite() bool {
	return l == 0
}
