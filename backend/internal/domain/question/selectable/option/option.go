package option

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	text      text.Text
	isCorrect bool
}

func New(t text.Text, isCorrect bool) (Option, error) {
	if err := validateOptionText(t); err != nil {
		return Option{}, err
	}

	id, err := uid.New()
	if err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		text:      t,
		isCorrect: isCorrect,
	}, nil
}

func Restore(id uuid.UUID, t text.Text, isCorrect bool) (Option, error) {
	if id == uuid.Nil {
		return Option{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if err := validateOptionText(t); err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		text:      t,
		isCorrect: isCorrect,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Text() text.Text {
	return o.text
}

func (o Option) IsCorrect() bool {
	return o.isCorrect
}
