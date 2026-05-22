package text

import "errors"

const (
	ValueCharsLimit int = 1000
)

var (
	ErrInvalid = errors.New("invalid value")
)

type Text struct {
	value string
}

func New(value string) (Text, error) {
	if err := validateValue(value); err != nil {
		return Text{}, err
	}

	return Text{
		value: value,
	}, nil
}

func (t Text) Value() string {
	return t.value
}

func (t Text) IsIncomplete() bool {
	return t.value == ""
}
