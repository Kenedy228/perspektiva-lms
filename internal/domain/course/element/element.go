package element

import (
	"gitflic.ru/lms/internal/domain/course/element/title"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Element struct {
	id      uuid.UUID
	t       title.Title
	content Content
}

func New(t title.Title, c Content) (*Element, error) {
	if err := validateContent(c); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Element{
		id:      id,
		t:       t,
		content: c.Clone(),
	}, nil
}

func (e *Element) ID() uuid.UUID {
	return e.id
}

func (e *Element) Title() title.Title {
	return e.t
}

func (e *Element) Content() Content {
	if e.content == nil {
		return nil
	}
	return e.content.Clone()
}

func (e *Element) ChangeTitle(t title.Title) {
	e.t = t
}

func (e *Element) ChangeContent(c Content) error {
	if err := validateContent(c); err != nil {
		return err
	}

	e.content = c.Clone()
	return nil
}

func (e *Element) Replicate() (*Element, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Element{
		id:      id,
		t:       e.t,
		content: e.content.Clone(),
	}, nil
}

func (e *Element) Clone() *Element {
	return &Element{
		id:      e.id,
		t:       e.t,
		content: e.content.Clone(),
	}
}

func (e *Element) IsInteractive() bool {
	if e.content == nil {
		return false
	}
	return e.content.IsInteractive()
}
