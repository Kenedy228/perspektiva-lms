//go:build legacy
// +build legacy

package element_test

import "gitflic.ru/lms/backend/internal/domain/question/content"

func makeContent(cType content.Type, value string) content.Content {
	c, _ := content.New(cType, value)
	return c
}
