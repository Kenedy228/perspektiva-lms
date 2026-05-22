package commands_test

import (
	"context"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	persondomain "gitflic.ru/lms/backend/internal/domain/person"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*persondomain.Person, error) {
	args := m.Called(ctx, id)

	var p *persondomain.Person
	if args.Get(0) != nil {
		p = args.Get(0).(*persondomain.Person)
	}

	return p, args.Error(1)
}

func (m *mockRepository) SNILSExists(ctx context.Context, value snils.SNILS, excludePersonID uuid.UUID) (bool, error) {
	args := m.Called(ctx, value, excludePersonID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, p *persondomain.Person) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockAuditRecorder struct {
	mock.Mock
}

func (m *mockAuditRecorder) RecordPersonAudit(ctx context.Context, event personports.AuditEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
