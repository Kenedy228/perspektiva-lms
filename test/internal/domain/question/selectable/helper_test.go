//go:build legacy
// +build legacy

package selectable_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/title"
)

func makeOptions(correctCount, incorrectCount int) []option.Option {
	opts := make([]option.Option, 0, correctCount+incorrectCount)

	for range correctCount {
		c := makeContent(content.TypeText, "foo")
		opt := makeOption(c, true)
		opts = append(opts, opt)
	}

	for range incorrectCount {
		c := makeContent(content.TypeText, "foo")
		opt := makeOption(c, false)
		opts = append(opts, opt)
	}

	return opts
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeOption(c content.Content, isCorrect bool) option.Option {
	opt, _ := option.New(c, isCorrect)
	return opt
}

func mockTitle() title.Title {
	t, _ := title.New(makeContent(content.TypeText, "text"))
	return t
}
