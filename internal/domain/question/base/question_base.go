package base

import (
	"time"

	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Base struct {
	id          uuid.UUID
	text        string
	description string
	image       uuid.UUID
	updatedAt   time.Time
	createdAt   time.Time
}

func New(params *Params) (Base, error) {
	if err := validateText(params.Text); err != nil {
		return Base{}, err
	}

	if err := validateDescription(params.Description); err != nil {
		return Base{}, err
	}

	id, err := utils.GenerateID()
	if err != nil {
		return Base{}, err
	}

	now := time.Now()

	return Base{
		id:          id,
		text:        params.Text,
		description: params.Description,
		image:       params.Image,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (b Base) ID() uuid.UUID {
	return b.id
}

func (b Base) Text() string {
	return b.text
}

func (b *Base) UpdateText(text string) error {
	if err := validateText(text); err != nil {
		return err
	}

	b.text = text
	b.updatedAt = time.Now()
	return nil
}

func (b Base) Description() string {
	return b.description
}

func (b Base) Image() uuid.UUID {
	return b.image
}

func (b *Base) UpdateImage(image uuid.UUID) {
	b.image = image
	b.updatedAt = time.Now()
}

func (b Base) HasImage() bool {
	return b.image != uuid.Nil
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
