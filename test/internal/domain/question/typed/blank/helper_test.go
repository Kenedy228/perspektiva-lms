package blank_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/title"
)

func makeContent(cType content.Type, value string) content.Content {
	c, _ := content.New(cType, value)
	return c
}

func makeContentSlice(count int) []content.Content {
	c := make([]content.Content, 0, count)

	for range count {
		c = append(c, makeContent(content.TypeText, "text"))
	}

	return c
}

func makeTitle(s string) title.Title {
	t, _ := title.New(makeContent(content.TypeText, s))
	return t
}
