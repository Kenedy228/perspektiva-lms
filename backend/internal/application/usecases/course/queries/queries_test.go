package queries

import (
	"context"
	"testing"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestListUsesStudentVisibilityForStudent(t *testing.T) {
	service := &queryService{views: []courseports.ShortView{{ID: uuid.New().String(), Title: "Course"}}}
	accountID := uuid.New()

	out, err := NewListQuery(service).Execute(context.Background(), ListInput{
		ActorRole: role.NewStudent(),
		AccountID: accountID.String(),
		Limit:     5,
	})
	if err != nil {
		t.Fatalf("list courses: %v", err)
	}
	if len(out.Views) != 1 {
		t.Fatalf("expected one course, got %d", len(out.Views))
	}
	if !service.studentCalled {
		t.Fatal("expected student visibility query")
	}
	if service.lastFilter.AccountID != accountID {
		t.Fatalf("expected account filter %s, got %s", accountID, service.lastFilter.AccountID)
	}
}

func TestRatingsRequiresManager(t *testing.T) {
	service := &queryService{}
	_, err := NewRatingsQuery(service).Execute(context.Background(), RatingsInput{
		ActorRole: role.NewStudent(),
		CourseID:  uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestStudentStatisticsRoleFilters(t *testing.T) {
	service := &queryService{}
	accountID := uuid.New()
	organizationID := uuid.New()

	_, err := NewStudentStatisticsQuery(service).Execute(context.Background(), StudentStatisticsInput{
		ActorRole: role.NewStudent(),
		AccountID: accountID.String(),
	})
	if err != nil {
		t.Fatalf("student statistics: %v", err)
	}
	if service.lastStatsFilter.AccountID != accountID {
		t.Fatalf("expected account filter %s, got %s", accountID, service.lastStatsFilter.AccountID)
	}

	_, err = NewStudentStatisticsQuery(service).Execute(context.Background(), StudentStatisticsInput{
		ActorRole:      role.NewOrganization(),
		OrganizationID: organizationID.String(),
	})
	if err != nil {
		t.Fatalf("organization statistics: %v", err)
	}
	if service.lastStatsFilter.OrganizationID != organizationID {
		t.Fatalf("expected organization filter %s, got %s", organizationID, service.lastStatsFilter.OrganizationID)
	}

	_, err = NewStudentStatisticsQuery(service).Execute(context.Background(), StudentStatisticsInput{
		ActorRole: role.NewCreator(),
	})
	if err == nil {
		t.Fatal("expected creator statistics to be forbidden")
	}
}

type queryService struct {
	lastFilter      courseports.Filter
	lastStatsFilter courseports.StudentStatisticsFilter
	studentCalled   bool
	views           []courseports.ShortView
}

func (s *queryService) ListManageable(_ context.Context, filter courseports.Filter) ([]courseports.ShortView, error) {
	s.lastFilter = filter
	return s.views, nil
}

func (s *queryService) ListVisibleForStudent(_ context.Context, filter courseports.Filter) ([]courseports.ShortView, error) {
	s.studentCalled = true
	s.lastFilter = filter
	return s.views, nil
}

func (s *queryService) GetDetailsByID(context.Context, uuid.UUID) (courseports.DetailedView, error) {
	return courseports.DetailedView{}, nil
}

func (s *queryService) ListRatings(context.Context, uuid.UUID, int, int) ([]courseports.StudentRatingView, error) {
	return nil, nil
}

func (s *queryService) ListStudentStatistics(_ context.Context, filter courseports.StudentStatisticsFilter) ([]courseports.StudentRatingView, error) {
	s.lastStatsFilter = filter
	return nil, nil
}
