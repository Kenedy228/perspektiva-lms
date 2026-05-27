package option

import (
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	value     string
	isCorrect bool
}

func New(value string, isCorrect bool) (Option, error) {
	if err := validateValue(value); err != nil {
		return Option{}, err
	}

	id, err := uid.New()
	if err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		value:     value,
		isCorrect: isCorrect,
	}, nil
}

func Restore(b base.Base, id uuid.UUID, value string, isCorrect bool) (Option, error) {
	if err := validateValue(value); err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		value:     t,
		isCorrect: isCorrect,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Text() text.Text {
	return o.value
}

func (o Option) IsCorrect() bool {
	return o.isCorrect
}
