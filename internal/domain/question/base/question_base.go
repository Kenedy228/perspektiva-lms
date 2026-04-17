package base

import (
	"time"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Base struct {
	id          uuid.UUID
	text        question.QText
	description question.QDescription
	imageID     uuid.UUID
	updatedAt   time.Time
	createdAt   time.Time
}

func New(params Params) (Base, error) {
	id, err := utils.GenerateID()
	if err != nil {
		return Base{}, err
	}

	now := time.Now()

	return Base{
		id:          id,
		text:        params.Text,
		description: params.Description,
		imageID:     params.ImageID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (b Base) ID() uuid.UUID {
	return b.id
}

func (b Base) Text() question.QText {
	return b.text
}

func (b Base) Description() question.QDescription {
	return b.description
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

func (b *Base) Touch() {
	b.updatedAt = time.Now()
}

func (b *Base) UpdateText(text question.QText) {
	if b.text == text {
		return
	}
	b.text = text
	b.updatedAt = time.Now()
}

func (b *Base) UpdateImage(imageID uuid.UUID) {
	if imageID == uuid.Nil {
		return
	}
	b.imageID = imageID
	b.updatedAt = time.Now()
}

func (b *Base) RemoveImage() {
	if b.imageID == uuid.Nil {
		return
	}

	b.imageID = uuid.Nil
	b.updatedAt = time.Now()
}

func (b Base) HasImage() bool {
	return b.imageID != uuid.Nil
}
