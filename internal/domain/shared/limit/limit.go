package limit

import "time"

type Limit struct {
	value int
}

func New(value int) (Limit, error) {
	if err := validateValue(value); err != nil {
		return Limit{}, err
	}

	return Limit{
		value: value,
	}, nil
}

func (l Limit) IsInfinite() bool {
	return l.value == 0
}

func (l Limit) Value() int {
	return l.value
}

func (l Limit) Duration() time.Duration {
	return time.Second * time.Duration(l.value)
}
