package base

import (
	"gitflic.ru/lms/internal/domain/question/attachment"
	"gitflic.ru/lms/internal/domain/question/title"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Base struct {
	id         uuid.UUID
	title      title.Title
	attachment *attachment.Attachment
}

func New(t title.Title) (*Base, error) {
	id, err := uid.New()
	if err != nil {
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

func (b *Base) Attachment() (attachment.Attachment, bool) {
	if b.attachment == nil {
		return attachment.Attachment{}, false
	}

	return *b.attachment, true
}

func (b *Base) ChangeTitle(t title.Title) {
	b.title = t
}

func (b *Base) HasAttachment() bool {
	if b.attachment == nil {
		return false
	}

	return true
}

func (b *Base) ChangeAttachment(a attachment.Attachment) {
	if b.HasAttachment() {
		b.RemoveAttachment()
	}

	b.attachment = &a
}

func (b *Base) RemoveAttachment() {
	if !b.HasAttachment() {
		return
	}

	b.attachment = nil
}

func (b *Base) Clone() *Base {
	var cAttachment *attachment.Attachment

	if b.attachment != nil {
		temp := *b.attachment
		cAttachment = &temp
	}

	return &Base{
		id:         b.id,
		title:      b.title,
		attachment: cAttachment,
	}
}
