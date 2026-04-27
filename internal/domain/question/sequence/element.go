package sequence

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Element struct {
	id      uuid.UUID
	content question.Content
}

func NewElement(params ElementParams) (Element, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return Element{}, err
	}

	return Element{
		id:      id,
		content: params.Content,
	}, nil
}

func (i Element) ID() uuid.UUID {
	return i.id
}

func (i Element) Content() question.Content {
	return i.content
}

func (i Element) Equal(other Element) bool {
	return i.content.Equal(other.content)
}
