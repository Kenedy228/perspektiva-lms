package commands_test

import (
	"context"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (question.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(question.Question), args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, q question.Question) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
