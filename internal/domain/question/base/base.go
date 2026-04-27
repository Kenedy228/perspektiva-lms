package base

import (
	"time"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Base struct {
	id        uuid.UUID
	text      string
	imageID   uuid.UUID
	updatedAt time.Time
	createdAt time.Time
}

func New(params Params) (*Base, error) {
	if err := validateText(params.Text); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Base{
		id:        id,
		text:      params.Text,
		imageID:   params.ImageID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (b Base) ID() uuid.UUID {
	return b.id
}

func (b Base) Text() string {
	return b.text
}

func (b Base) ImageID() uuid.UUID {
	return b.imageID
}

func (b Base) CreatedAt() time.Time {
	return b.createdAt
}

func (b Base) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b Base) HasImage() bool {
	return b.imageID != uuid.Nil
}

func (b *Base) RemoveImage() {
	if b.imageID == uuid.Nil {
		return
	}

	b.imageID = uuid.Nil
	b.Touch()
}

func (b *Base) Touch() {
	b.updatedAt = time.Now()
}

func (b *Base) UpdateImage(imageID uuid.UUID) {
	if imageID == uuid.Nil {
		b.RemoveImage()
		return
	}

	b.imageID = imageID
	b.Touch()
}

func (b *Base) UpdateText(text string) error {
	if err := validateText(text); err != nil {
		return err
	}

	if b.text == text {
		return nil
	}

	b.text = text
	b.Touch()
	return nil
}

func (b *Base) Clone() *Base {
	return &Base{
		id:        b.id,
		text:      b.text,
		imageID:   b.imageID,
		createdAt: b.createdAt,
		updatedAt: b.updatedAt,
	}
}
