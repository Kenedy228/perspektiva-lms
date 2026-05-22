package variant

import (
	"errors"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

const TextCharsLimit int = 1000

var ErrInvalid = errors.New("invalid value")

type Variant struct {
	text text.Text
}

func New(t text.Text) (Variant, error) {
	if err := validateText(t); err != nil {
		return Variant{}, err
	}

	return Variant{
		text: t,
	}, nil
}

func (v Variant) Text() text.Text {
	return v.text
}

func (v Variant) IsZero() bool {
	return len(v.text.Value()) == 0
}
