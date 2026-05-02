package answer_test

import "gitflic.ru/lms/internal/domain/question"

type mockAnswer struct{}

func (m mockAnswer) Clone() question.Answer {
	return m
}

func (m mockAnswer) IsEmpty() bool {
	panic("")
}
