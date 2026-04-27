package block_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/block"
	"gitflic.ru/lms/internal/domain/element"
	"github.com/stretchr/testify/assert"
)

type blockBuilder struct {
	title    string
	elements []*element.Block
}

func newBlockBuilder() *blockBuilder {
	return &blockBuilder{}
}

func (b *blockBuilder) withTitle(title string) *blockBuilder {
	b.title = title
	return b
}

func (b *blockBuilder) withElement() *blockBuilder {
	b.elements = append(b.elements, newElement())
	return b
}

func (b *blockBuilder) build(t *testing.T, wantErr error) *block.Block {
	params := block.Params{
		Title:    b.title,
		Elements: b.elements,
	}

	bl, err := block.New(params)
	assert.ErrorIs(t, err, wantErr)

	return bl
}
