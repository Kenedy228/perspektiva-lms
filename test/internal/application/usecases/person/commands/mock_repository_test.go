package commands_test

import (
	"context"

	"gitflic.ru/lms/internal/domain/person"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*person.Person, error) {
	args := m.Called(ctx, id)

	var p *person.Person
	if args.Get(0) != nil {
		p = args.Get(0).(*person.Person)
	}

	return p, args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, p *person.Person) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
