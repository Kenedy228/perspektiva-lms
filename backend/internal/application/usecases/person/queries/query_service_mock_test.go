package queries_test

import (
	"context"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockQueryService struct {
	mock.Mock
}

func (m *mockQueryService) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, organizationID, limit, offset)

	var views []personports.PersonShortView
	if args.Get(0) != nil {
		views = args.Get(0).([]personports.PersonShortView)
	}

	return views, args.Error(1)
}

func (m *mockQueryService) ListByLastName(ctx context.Context, lastName string, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, lastName, limit, offset)

	var views []personports.PersonShortView
	if args.Get(0) != nil {
		views = args.Get(0).([]personports.PersonShortView)
	}

	return views, args.Error(1)
}

func (m *mockQueryService) ListBySnils(ctx context.Context, snils string, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, snils, limit, offset)

	var views []personports.PersonShortView
	if args.Get(0) != nil {
		views = args.Get(0).([]personports.PersonShortView)
	}

	return views, args.Error(1)
}

func (m *mockQueryService) GetDetailsByID(ctx context.Context, id uuid.UUID) (personports.PersonDetailedView, error) {
	args := m.Called(ctx, id)

	var view personports.PersonDetailedView
	if args.Get(0) != nil {
		view = args.Get(0).(personports.PersonDetailedView)
	}

	return view, args.Error(1)
}
