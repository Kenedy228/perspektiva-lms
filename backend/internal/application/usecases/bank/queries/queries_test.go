package queries

import (
	"context"
	"errors"
	"testing"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestListNormalizesFilters(t *testing.T) {
	ctx := context.Background()
	questionID := uuid.New()
	service := &queryService{
		views: []bankports.ShortView{{ID: uuid.New().String(), Title: "Bank", QuestionsCount: 3}},
	}

	out, err := NewListQuery(service).Execute(ctx, ListInput{
		ActorRole:     role.NewCreator(),
		TitleContains: "  final   bank  ",
		QuestionID:    questionID.String(),
		MinQuestions:  1,
		MaxQuestions:  5,
		Limit:         10,
		Offset:        2,
	})
	if err != nil {
		t.Fatalf("list banks: %v", err)
	}
	if len(out.Views) != 1 {
		t.Fatalf("expected one view, got %d", len(out.Views))
	}
	if service.lastFilter.TitleContains != "final bank" {
		t.Fatalf("expected normalized title filter, got %q", service.lastFilter.TitleContains)
	}
	if service.lastFilter.QuestionID != questionID {
		t.Fatalf("expected question filter %s, got %s", questionID, service.lastFilter.QuestionID)
	}
	if service.lastFilter.Limit != 10 || service.lastFilter.Offset != 2 {
		t.Fatalf("unexpected pagination: %+v", service.lastFilter)
	}
}

func TestListRejectsForbiddenAndInvalidFilters(t *testing.T) {
	ctx := context.Background()
	service := &queryService{}

	_, err := NewListQuery(service).Execute(ctx, ListInput{
		ActorRole: role.NewStudent(),
	})
	if !errors.Is(err, common.ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}

	_, err = NewListQuery(service).Execute(ctx, ListInput{
		ActorRole:    role.NewAdmin(),
		MinQuestions: 10,
		MaxQuestions: 1,
	})
	if !errors.Is(err, common.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestGetDetailsByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()
	service := &queryService{
		details: bankports.DetailedView{ID: id.String(), Title: "Bank"},
	}

	view, err := NewGetDetailsByIDQuery(service).Execute(ctx, GetDetailsByIDInput{
		ActorRole: role.NewAdmin(),
		BankID:    id.String(),
	})
	if err != nil {
		t.Fatalf("get details: %v", err)
	}
	if view.ID != id.String() {
		t.Fatalf("expected view id %s, got %s", id, view.ID)
	}
}

type queryService struct {
	lastFilter bankports.Filter
	views      []bankports.ShortView
	details    bankports.DetailedView
}

func (s *queryService) List(_ context.Context, filter bankports.Filter) ([]bankports.ShortView, error) {
	s.lastFilter = filter
	return s.views, nil
}

func (s *queryService) GetDetailsByID(_ context.Context, _ uuid.UUID) (bankports.DetailedView, error) {
	return s.details, nil
}
