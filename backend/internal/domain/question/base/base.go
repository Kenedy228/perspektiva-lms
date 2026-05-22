package base

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Base struct {
	id            uuid.UUID
	title         title.Title
	attachment    attachment.Attachment
	hasAttachment bool
}

func New(t title.Title) (*Base, error) {
	if t.IsZero() {
		return nil, fmt.Errorf("%w: invalid value", ErrInvalid)
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

func Restore(id uuid.UUID, t title.Title, att *attachment.Attachment) (*Base, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if t.IsZero() {
		return nil, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	b := &Base{
		id:    id,
		title: t,
	}

	if att != nil && !att.IsZero() {
		b.attachment = *att
		b.hasAttachment = true
	}

	return b, nil
}

func (b *Base) ID() uuid.UUID {
	return b.id
}

func (b *Base) Title() title.Title {
	return b.title
}

func (b *Base) Attachment() (attachment.Attachment, bool) {
	if !b.hasAttachment {
		return attachment.Attachment{}, false
	}

	return b.attachment, true
}

func (b *Base) ChangeTitle(t title.Title) {
	if t.IsZero() {
		return
	}

	b.title = t
}

func (b *Base) ChangeAttachment(a attachment.Attachment) {
	if a.IsZero() {
		return
	}

	b.attachment = a
	b.hasAttachment = true
}

func (b *Base) HasAttachment() bool {
	return b.hasAttachment
}

func (b *Base) RemoveAttachment() {
	b.attachment = attachment.Attachment{}
	b.hasAttachment = false
}

func (b *Base) Clone() *Base {
	return &Base{
		id:            b.id,
		title:         b.title,
		attachment:    b.attachment,
		hasAttachment: b.hasAttachment,
	}
}
