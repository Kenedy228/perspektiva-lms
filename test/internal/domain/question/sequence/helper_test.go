package sequence_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/sequence/option"
	"gitflic.ru/lms/internal/domain/question/title"
)

func makeOptions(count int) []option.Option {
	opts := make([]option.Option, 0, count)

	for range count {
		c := makeContent(content.TypeText, "foo")
		opt := makeOption(c)
		opts = append(opts, opt)
	}

	return opts
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeOption(c content.Content) option.Option {
	opt, _ := option.New(c)
	return opt
}

func mockTitle() title.Title {
	t, _ := title.New(makeContent(content.TypeText, "text"))
	return t
}
