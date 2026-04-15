package sequence

import (
	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Item struct {
	id      uuid.UUID
	content content.RichContent
}

func NewItem(content content.RichContent) (Item, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return Item{}, err
	}

	return Item{
		id:      id,
		content: content,
	}, nil
}

func (i Item) ID() uuid.UUID {
	return i.id
}

func (i Item) Content() content.RichContent {
	return i.content
}
