package block_test

import (
	"gitflic.ru/lms/internal/domain/element"
)

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

func newElement() *element.Block {
	params := element.Params{
		Title:   "title",
		Content: &contentFixture{},
	}

	el, _ := element.New(params)
	return el
}
