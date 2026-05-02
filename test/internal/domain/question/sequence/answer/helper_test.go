package answer_test

import (
	"gitflic.ru/lms/internal/domain/question/sequence/answer"
	"github.com/google/uuid"
)

func makeAnswerOptions(count int) []answer.AnswerOption {
	opts := make([]answer.AnswerOption, 0, count)

	for range count {
		opts = append(opts, answer.AnswerOption{
			OptionID: uuid.New(),
		})
	}

	return opts
}
