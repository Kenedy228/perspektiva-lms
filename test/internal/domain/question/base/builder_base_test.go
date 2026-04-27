package base_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/base"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type baseBuilder struct {
	text    string
	imageID uuid.UUID
}

func newBaseBuilder() *baseBuilder {
	return &baseBuilder{
		text:    "",
		imageID: uuid.Nil,
	}
}

func (b *baseBuilder) withText(s string) *baseBuilder {
	b.text = s
	return b
}

func (b *baseBuilder) withImage(id uuid.UUID) *baseBuilder {
	b.imageID = id
	return b
}

func (b *baseBuilder) build(t *testing.T, wantErr error) *base.Base {
	t.Helper()
	params := base.Params{
		Text:    b.text,
		ImageID: b.imageID,
	}

	base, err := base.New(params)
	assert.ErrorIs(t, err, wantErr)

	return base
}
