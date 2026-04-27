package element

import (
	"time"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Element struct {
	id        uuid.UUID
	title     string
	content   Content
	createdAt time.Time
	updatedAt time.Time
}

func New(params Params) (*Element, error) {
	title := normalizeTitle(params.Title)

	if err := validateTitle(title); err != nil {
		return nil, err
	}

	if err := validateContent(params.Content); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Element{
		id:        id,
		title:     title,
		content:   params.Content.Clone(),
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (e *Element) ID() uuid.UUID {
	return e.id
}

func (e *Element) Title() string {
	return e.title
}

func (e *Element) Content() Content {
	if e.content == nil {
		return nil
	}
	return e.content.Clone()
}

func (e *Element) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Element) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Element) ChangeTitle(title string) error {
	title = normalizeTitle(title)
	if err := validateTitle(title); err != nil {
		return err
	}

	e.title = title
	e.updatedAt = time.Now()
	return nil
}

func (e *Element) ChangeContent(content Content) error {
	if err := validateContent(content); err != nil {
		return err
	}

	e.content = content.Clone()
	e.updatedAt = time.Now()
	return nil
}

func (e *Element) Copy() (*Element, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Element{
		id:        id,
		title:     e.title,
		content:   e.content.Clone(),
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (e *Element) Clone() *Element {
	return &Element{
		id:        e.id,
		title:     e.title,
		content:   e.content.Clone(),
		createdAt: e.createdAt,
		updatedAt: e.updatedAt,
	}
}

func (e *Element) IsInteractive() bool {
	if e.content == nil {
		return false
	}
	return e.content.IsInteractive()
}
