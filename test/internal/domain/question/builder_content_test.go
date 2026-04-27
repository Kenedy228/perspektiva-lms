package question_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
)

type contentBuilder struct {
	cType question.ContentType
	value string
}

func newContentBuilder() *contentBuilder {
	return &contentBuilder{
		cType: question.ContentType(""),
		value: "",
	}
}

func (b *contentBuilder) withTextCType() *contentBuilder {
	b.cType = question.ContentTypeText
	return b
}

func (b *contentBuilder) withAudioCType() *contentBuilder {
	b.cType = question.ContentTypeAudio
	return b
}

func (b *contentBuilder) withImageCType() *contentBuilder {
	b.cType = question.ContentTypeImage
	return b
}

func (b *contentBuilder) withCType(val string) *contentBuilder {
	b.cType = question.ContentType(val)
	return b
}

func (b *contentBuilder) withValue(val string) *contentBuilder {
	b.value = val
	return b
}

func (b *contentBuilder) build(t *testing.T, wantErr error) question.Content {
	t.Helper()
	c, err := question.NewContent(b.cType, b.value)
	assert.ErrorIs(t, err, wantErr)

	return c
}
