//go:build legacy
// +build legacy

package quiz_test

import (
	limit2 "gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
)

func makeTitle(val string) title.Title {
	t, err := title.New(val)
	if err != nil {
		panic(err)
	}
	return t
}

func makeTime(val int) limit2.Time {
	l, err := limit2.NewTime(val)
	if err != nil {
		panic(err)
	}
	return l
}

func makeAttempts(val int) limit2.Attempts {
	l, err := limit2.NewAttempts(val)
	if err != nil {
		panic(err)
	}
	return l
}
