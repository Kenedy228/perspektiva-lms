package grading

import (
	"gitflic.ru/lms/internal/domain/grading/score"
	"gitflic.ru/lms/internal/domain/question"
)

type Checker interface {
	Check(q question.Question, a question.Answer) (score.Score, error)
	Supports(t question.Type) bool
}
