//go:build legacy
// +build legacy

package answer_test

import (
	"gitflic.ru/lms/backend/internal/domain/question"
)

type mockAnswer struct{}

func (m mockAnswer) Clone() question.Answer {
	return m
}

func (m mockAnswer) IsEmpty() bool {
	panic("")
}
