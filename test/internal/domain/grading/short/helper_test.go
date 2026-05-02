package short_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/short/answer"
	"gitflic.ru/lms/internal/domain/question/short/variant"
	"gitflic.ru/lms/internal/domain/question/title"
)

func makeTitle() title.Title {
	t, _ := title.New(makeContent(content.TypeText, "title"))
	return t
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func makeVariant(s string) variant.Variant {
	v, _ := variant.New(makeContent(content.TypeText, s))
	return v
}

func makeAnswerVariant(s string) answer.AnswerVariant {
	return answer.AnswerVariant{
		Input: s,
	}
}
