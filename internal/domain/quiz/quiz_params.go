package quiz

import "gitflic.ru/lms/internal/domain/shared/limit"

type Params struct {
	Title       string
	Sources     []Source
	MaxAttempts int
	TimeLimit   limit.Limit
}
