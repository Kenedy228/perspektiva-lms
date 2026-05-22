//go:build legacy
// +build legacy

package typed_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/question/content"
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
