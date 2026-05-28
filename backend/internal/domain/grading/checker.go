package grading

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
)

// Checker описывает сервис проверки ответа для конкретного типа вопроса.
type Checker interface {
	// Check вычисляет итоговый score за ответ на вопрос.
	Check(q question.Question, a question.Answer) (score.Score, error)
}
