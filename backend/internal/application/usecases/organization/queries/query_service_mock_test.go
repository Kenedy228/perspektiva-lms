package queries_test

import (
	"context"

	"gitflic.ru/lms/backend/internal/application/ports/organization"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockQueryService struct {
	mock.Mock
}

func (s *mockQueryService) ListByName(ctx context.Context, name string, limit, offset int) ([]organization.OrganizationShortView, error) {
	args := s.Called(ctx, name, limit, offset)

	var views []organization.OrganizationShortView
	if args.Get(0) != nil {
		views = args.Get(0).([]organization.OrganizationShortView)
	}

	return views, args.Error(1)
}

func (s *mockQueryService) ListByINN(ctx context.Context, inn string, limit, offset int) ([]organization.OrganizationShortView, error) {
	args := s.Called(ctx, inn, limit, offset)

	var views []organization.OrganizationShortView
	if args.Get(0) != nil {
		views = args.Get(0).([]organization.OrganizationShortView)
	}

	return views, args.Error(1)
}

func (s *mockQueryService) GetDetailsByID(ctx context.Context, id uuid.UUID) (organization.OrganizationDetailedView, error) {
	args := s.Called(ctx, id)

	var view organization.OrganizationDetailedView
	if args.Get(0) != nil {
		view = args.Get(0).(organization.OrganizationDetailedView)
	}

	return view, args.Error(1)
}
