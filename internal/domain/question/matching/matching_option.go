package matching

import (
	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Option struct {
	id      uuid.UUID
	content content.RichContent
}

func NewOption(content content.RichContent) (Option, error) {
	id, err := utils.GenerateID()

	if err != nil {
		return Option{}, err
	}

	return Option{
		id:      id,
		content: content,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Content() content.RichContent {
	return o.content
}
