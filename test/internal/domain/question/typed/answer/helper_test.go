package answer_test

import "gitflic.ru/lms/internal/domain/question/typed/answer"

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
