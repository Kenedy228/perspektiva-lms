//go:build legacy
// +build legacy

package answer_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/typed/answer"
)

func makeAnswerBlanks(count int) []answer.AnswerBlank {
	blanks := make([]answer.AnswerBlank, 0, count)

	for range count {
		blanks = append(blanks, answer.AnswerBlank{
			Placeholder: "placeholder",
			Variant:     "variant",
		})
	}

	return blanks
}
