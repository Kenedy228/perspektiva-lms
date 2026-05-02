package typed_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/typed/blank"
)

func makeContent(value string) content.Content {
	c, _ := content.New(content.TypeText, value)
	return c
}

func makeBlank(placeholder, variantVal string) blank.Blank {
	c := makeContent(variantVal)
	b, _ := blank.New(placeholder, []content.Content{c})
	return b
}
