//go:build legacy
// +build legacy

package answer_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"github.com/google/uuid"
)

func makeOptions(size int) []answer.AnswerOption {
	answerOptions := make([]answer.AnswerOption, 0, size)

	for range size {
		answerOptions = append(answerOptions, answer.AnswerOption{
			OptionID: uuid.New(),
		})
	}

	return answerOptions
}
