package quiz

import "gitflic.ru/lms/internal/domain/quiz/source"

type Params struct {
	Title        string
	Sources      []source.Source
	AttemptLimit int
	TimeLimit    int
}
