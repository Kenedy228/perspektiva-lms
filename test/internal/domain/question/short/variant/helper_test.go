//go:build legacy
// +build legacy

package variant_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"gitflic.ru/lms/backend/internal/domain/question/content"
)

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeVariant(c content.Content) (variant.Variant, error) {
	return variant.New(c)
}
