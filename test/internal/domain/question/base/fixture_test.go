//go:build legacy
// +build legacy

package base_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/title"
)

func makeTitle(s string) title.Title {
	t, _ := title.New(makeContent(content.TypeText, s))
	return t
}

func makeAttachment(s string) attachment.Attachment {
	a, _ := attachment.New(makeContent(content.TypeImage, s))
	return a
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(content.TypeText, s)
	return c
}
