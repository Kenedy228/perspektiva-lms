package commands_test

import (
	"context"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/domain/organization"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error) {
	args := m.Called(ctx, id)

	var org *organization.Organization
	if args.Get(0) != nil {
		org = args.Get(0).(*organization.Organization)
	}

	return org, args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, o *organization.Organization) error {
	args := m.Called(ctx, o)

	return args.Error(0)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

type mockAuditRecorder struct {
	mock.Mock
}

func (m *mockAuditRecorder) RecordOrganizationAudit(ctx context.Context, event orgports.AuditEvent) error {
	args := m.Called(ctx, event)

	return args.Error(0)
}
