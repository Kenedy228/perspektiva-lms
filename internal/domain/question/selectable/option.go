package selectable

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	content   question.Content
	isCorrect bool
}

func NewItem(params OptionParams) (Option, error) {
	id, err := uid.New()
	if err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		content:   params.Content,
		isCorrect: params.IsCorrect,
	}, nil
}

func (i Option) ID() uuid.UUID {
	return i.id
}

func (i Option) Content() question.Content {
	return i.content
}

func (i Option) IsCorrect() bool {
	return i.isCorrect
}

func (i Option) Equal(other Option) bool {
	return i.content == other.content
}
