package quiz_test

import (
	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/stretchr/testify/mock"
)

type mockCriteria struct {
	mock.Mock
}

func (m *mockCriteria) QuestionCount() int {
	panic("")
}

func (m *mockCriteria) Type() criteria.Type {
	panic("")
}
