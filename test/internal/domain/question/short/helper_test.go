//go:build legacy
// +build legacy

package short_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/title"
)

func makeVariants(count int) []variant.Variant {
	variants := make([]variant.Variant, 0, count)

	for range count {
		variants = append(
			variants,
			makeVariant(makeContent(content.TypeText, "value")),
		)
	}

	return variants
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeVariant(c content.Content) variant.Variant {
	v, _ := variant.New(c)
	return v
}

func makeTitle() title.Title {
	t, _ := title.New(makeContent(content.TypeText, "text"))
	return t
}
