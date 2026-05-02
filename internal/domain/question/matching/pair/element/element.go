package element

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Element struct {
	id uuid.UUID
	c  content.Content
}

func New(c content.Content) (Element, error) {
	id, err := uid.New()
	if err != nil {
		return Element{}, err
	}

	return Element{
		id: id,
		c:  c,
	}, nil
}

func (e Element) ID() uuid.UUID {
	return e.id
}

func (e Element) Content() content.Content {
	return e.c
}

func (e Element) Value() string {
	return e.c.Value()
}
