package attempt_test

import (
	"time"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockQuestion struct {
	mock.Mock
}

func (m *mockQuestion) ID() uuid.UUID {
	panic("")
}

func (m *mockQuestion) Clone() question.Question {
	args := m.Called()
	return args.Get(0).(question.Question)
}

func (m *mockQuestion) CheckAnswer(answer question.Answer) bool {
	args := m.Called(answer)
	return args.Bool(0)
}

func (m *mockQuestion) Text() string {
	panic("")
}

func (m *mockQuestion) Instruction() string {
	panic("")
}

func (m *mockQuestion) ImageID() uuid.UUID {
	panic("")
}

func (m *mockQuestion) CreatedAt() time.Time {
	panic("")
}

func (m *mockQuestion) UpdatedAt() time.Time {
	panic("")
}

func (m *mockQuestion) HasImage() bool {
	panic("")
}

func (m *mockQuestion) Type() question.Type {
	panic("")
}

type mockAnswer struct {
	mock.Mock
}

func (m *mockAnswer) IsEmpty() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *mockAnswer) Clone() question.Answer {
	args := m.Called()
	return args.Get(0).(question.Answer)
}
