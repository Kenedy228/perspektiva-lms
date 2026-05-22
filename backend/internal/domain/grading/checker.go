package grading

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
)

type Checker interface {
	Check(q question2.Question, a question2.Answer) (score.Score, error)
	Supports(t question2.Type) bool
}
