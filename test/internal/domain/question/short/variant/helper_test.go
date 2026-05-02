package variant_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/short/variant"
)

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeVariant(c content.Content) (variant.Variant, error) {
	return variant.New(c)
}
