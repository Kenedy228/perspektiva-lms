//go:build legacy
// +build legacy

package element_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/element"
	"github.com/stretchr/testify/assert"
)

type elementBuilder struct {
	title   string
	content element.Content
}

type contentFixture struct{}

func (f *contentFixture) Type() element.Type {
	panic("")
}

func (f *contentFixture) IsInteractive() bool {
	panic("")
}

func (f *contentFixture) Clone() element.Content {
	return f
}

func newElementBuilder() *elementBuilder {
	return &elementBuilder{}
}

func (b *elementBuilder) withTitle(title string) *elementBuilder {
	b.title = title
	return b
}

func (b *elementBuilder) withContent() *elementBuilder {
	b.content = &contentFixture{}
	return b
}

func (b *elementBuilder) build(t *testing.T, wantErr error) *element.Element {
	params := element.Params{
		Title:   b.title,
		Content: b.content,
	}

	e, err := element.New(params)
	assert.ErrorIs(t, err, wantErr)

	return e
}
