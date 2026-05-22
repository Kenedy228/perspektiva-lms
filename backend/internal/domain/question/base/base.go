package base

import (
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Base struct {
	id    uuid.UUID
	title title.Title
}

func New(t title.Title) (*Base, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Base{
		id:    id,
		title: t,
	}, nil
}

func Restore(id uuid.UUID, t title.Title) (*Base, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateTitle(t); err != nil {
		return nil, err
	}

	return &Base{
		id:    id,
		title: t,
	}, nil
}

func (b *Base) ID() uuid.UUID {
	return b.id
}

func (b *Base) Title() title.Title {
	return b.title
}

func (b *Base) ChangeTitle(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}

	b.title = t
	return nil
}

func (b *Base) Clone() *Base {
	return &Base{
		id:    b.id,
		title: b.title,
	}
}
