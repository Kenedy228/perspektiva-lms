package grading

import "gitflic.ru/lms/backend/internal/domain/question"

type AnswerValidator interface {
	Validate(q question.Question, a question.Answer) error
}
