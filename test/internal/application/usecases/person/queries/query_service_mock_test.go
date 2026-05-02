package queries_test

import (
	"context"

	personports "gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockQueryService struct {
	mock.Mock
}

func (m *mockQueryService) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, organizationID, limit, offset)

	var persons []personports.PersonShortView
	if args.Get(0) != nil {
		persons = args.Get(0).([]personports.PersonShortView)
	}

	return persons, args.Error(1)
}

func (m *mockQueryService) ListByLastName(ctx context.Context, lastName string, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, lastName, limit, offset)

	var persons []personports.PersonShortView
	if args.Get(0) != nil {
		persons = args.Get(0).([]personports.PersonShortView)
	}

	return persons, args.Error(1)
}

func (m *mockQueryService) ListBySnils(ctx context.Context, snils string, limit, offset int) ([]personports.PersonShortView, error) {
	args := m.Called(ctx, snils, limit, offset)

	var persons []personports.PersonShortView
	if args.Get(0) != nil {
		persons = args.Get(0).([]personports.PersonShortView)
	}

	return persons, args.Error(1)
}

func (m *mockQueryService) GetDetailsByID(ctx context.Context, id uuid.UUID) (personports.PersonDetailedView, error) {
	args := m.Called(ctx, id)

	var person personports.PersonDetailedView
	if args.Get(0) != nil {
		person = args.Get(0).(personports.PersonDetailedView)
	}

	return person, args.Error(1)
}
