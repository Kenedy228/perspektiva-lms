package quiz_test

import (
	"gitflic.ru/lms/internal/domain/quiz/limit"
	"gitflic.ru/lms/internal/domain/quiz/title"
)

func makeTitle(val string) title.Title {
	t, err := title.New(val)
	if err != nil {
		panic(err)
	}
	return t
}

func makeTime(val int) limit.Time {
	l, err := limit.NewTime(val)
	if err != nil {
		panic(err)
	}
	return l
}

func makeAttempts(val int) limit.Attempts {
	l, err := limit.NewAttempts(val)
	if err != nil {
		panic(err)
	}
	return l
}
