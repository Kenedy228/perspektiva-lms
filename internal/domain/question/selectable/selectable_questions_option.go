package selectable

import (
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	text      string
	isCorrect bool
}

func NewOption(text string, isCorrect bool) (Option, error) {
	if err := validateOptionText(text); err != nil {
		return Option{}, err
	}

	id, err := utils.GenerateID()
	if err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		text:      text,
		isCorrect: isCorrect,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Text() string {
	return o.text
}

func (o Option) IsCorrect() bool {
	return o.isCorrect
}
