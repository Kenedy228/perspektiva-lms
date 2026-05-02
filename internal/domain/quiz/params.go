package quiz

import (
	"gitflic.ru/lms/internal/domain/quiz/limit"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"gitflic.ru/lms/internal/domain/quiz/title"
)

type Params struct {
	Title       title.Title
	MaxAttempts limit.Attempts
	TimeLimit   limit.Time
	Sources     []source.Source
}
