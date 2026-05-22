package queries_test

import (
	"context"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockQueryService struct {
	mock.Mock
}

func (m *mockQueryService) GetByID(ctx context.Context, id uuid.UUID) (accountports.AccountView, error) {
	args := m.Called(ctx, id)

	var view accountports.AccountView
	if args.Get(0) != nil {
		view = args.Get(0).(accountports.AccountView)
	}

	return view, args.Error(1)
}

func (m *mockQueryService) GetByPersonID(ctx context.Context, personID uuid.UUID) (accountports.AccountView, error) {
	args := m.Called(ctx, personID)

	var view accountports.AccountView
	if args.Get(0) != nil {
		view = args.Get(0).(accountports.AccountView)
	}

	return view, args.Error(1)
}

func (m *mockQueryService) List(ctx context.Context, filter accountports.ListFilter, limit, offset int) ([]accountports.AccountView, error) {
	args := m.Called(ctx, filter, limit, offset)

	var views []accountports.AccountView
	if args.Get(0) != nil {
		views = args.Get(0).([]accountports.AccountView)
	}

	return views, args.Error(1)
}
