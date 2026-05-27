package option

import (
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	value     string
	isCorrect bool
}

func New(value string, isCorrect bool) (Option, error) {
	value = normalizeValue(value)

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

func Restore(id uuid.UUID, value string, isCorrect bool) (Option, error) {
	if err := validateID(id); err != nil {
		return Option{}, err
	}

	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		value:     value,
		isCorrect: isCorrect,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Value() string {
	return o.value
}

func (o Option) IsCorrect() bool {
	return o.isCorrect
}

func (o Option) IsZero() bool {
	return o.id == uuid.Nil || o.value == ""
}
