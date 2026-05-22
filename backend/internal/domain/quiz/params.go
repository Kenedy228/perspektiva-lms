package quiz

import (
	limit2 "gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
)

type Params struct {
	Title            title.Title
	MaxAttempts      limit2.Attempts
	TimeLimit        limit2.Time
	ShuffleQuestions bool
	Sources          []source.Source
}
