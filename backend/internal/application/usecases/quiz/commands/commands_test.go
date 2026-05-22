package commands

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestCreateAllowsAdminAndStoresQuiz(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	bankID := uuid.New()
	inspector := &bankInspector{
		counts: map[uuid.UUID]int{bankID: 5},
	}

	out, err := NewCreateUseCase(repo, inspector).Execute(ctx, CreateInput{
		ActorRole:        role.NewAdmin(),
		Title:            "Final test",
		MaxAttempts:      2,
		TimeLimitSeconds: 900,
		ShuffleQuestions: true,
		Sources: []SourceInput{{
			BankID:        bankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 5,
		}},
	})
	if err != nil {
		t.Fatalf("create quiz: %v", err)
	}
	if out.ID == "" {
		t.Fatal("expected output id")
	}

	saved := repo.saved
	if saved == nil {
		t.Fatal("expected quiz to be saved")
	}
	if !saved.ShuffleQuestions() {
		t.Fatal("expected quiz shuffle policy to be enabled")
	}
	if saved.Attempts().Count() != 2 {
		t.Fatalf("expected attempts limit 2, got %d", saved.Attempts().Count())
	}
	if saved.Time().Seconds() != 900 {
		t.Fatalf("expected time limit 900, got %d", saved.Time().Seconds())
	}
}

func TestCreateAllowsCreator(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	bankID := uuid.New()
	inspector := &bankInspector{counts: map[uuid.UUID]int{bankID: 1}}

	_, err := NewCreateUseCase(repo, inspector).Execute(ctx, CreateInput{
		ActorRole: role.NewCreator(),
		Title:     "Creator test",
		Sources: []SourceInput{{
			BankID:        bankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 1,
		}},
	})
	if err != nil {
		t.Fatalf("creator should create quiz: %v", err)
	}
}

func TestCreateRejectsStudent(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	bankID := uuid.New()
	inspector := &bankInspector{counts: map[uuid.UUID]int{bankID: 1}}

	_, err := NewCreateUseCase(repo, inspector).Execute(ctx, CreateInput{
		ActorRole: role.NewStudent(),
		Title:     "Student test",
		Sources: []SourceInput{{
			BankID:        bankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 1,
		}},
	})
	if !errors.Is(err, common.ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
	if repo.saved != nil {
		t.Fatal("student-created quiz must not be saved")
	}
}

func TestCreateRejectsSourceCountGreaterThanBankCount(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	bankID := uuid.New()
	inspector := &bankInspector{counts: map[uuid.UUID]int{bankID: 2}}

	_, err := NewCreateUseCase(repo, inspector).Execute(ctx, CreateInput{
		ActorRole: role.NewAdmin(),
		Title:     "Oversized test",
		Sources: []SourceInput{{
			BankID:        bankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 3,
		}},
	})
	if !errors.Is(err, common.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
	if repo.saved != nil {
		t.Fatal("invalid quiz must not be saved")
	}
}

func TestCreateRejectsManualQuestionsOutsideBank(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	bankID := uuid.New()
	firstQuestionID := uuid.New()
	inspector := &bankInspector{
		counts:       map[uuid.UUID]int{bankID: 1},
		belongsToAll: false,
	}

	_, err := NewCreateUseCase(repo, inspector).Execute(ctx, CreateInput{
		ActorRole: role.NewAdmin(),
		Title:     "Manual test",
		Sources: []SourceInput{{
			BankID:       bankID.String(),
			CriteriaType: criteria.TypeManual.String(),
			QuestionIDs:  []string{firstQuestionID.String()},
		}},
	})
	if !errors.Is(err, common.ErrInvalidInput) {
		t.Fatalf("expected invalid input error, got %v", err)
	}
}

func TestReplaceSourcesChecksNewSourceCapacity(t *testing.T) {
	ctx := context.Background()
	repo := newQuizRepo()
	originalBankID := uuid.New()
	replacementBankID := uuid.New()
	inspector := &bankInspector{counts: map[uuid.UUID]int{
		originalBankID:    1,
		replacementBankID: 4,
	}}

	q := mustStoredQuiz(t, repo, originalBankID)

	err := NewReplaceSourcesUseCase(repo, inspector).Execute(ctx, ReplaceSourcesInput{
		ActorRole: role.NewCreator(),
		QuizID:    q.ID().String(),
		Sources: []SourceInput{{
			BankID:        replacementBankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 4,
		}},
	})
	if err != nil {
		t.Fatalf("replace sources: %v", err)
	}
	if repo.saved.Sources()[0].BankID() != replacementBankID {
		t.Fatal("expected replacement source to be saved")
	}
}

type quizRepo struct {
	items map[uuid.UUID]*quizdomain.Quiz
	saved *quizdomain.Quiz
}

func newQuizRepo() *quizRepo {
	return &quizRepo{items: make(map[uuid.UUID]*quizdomain.Quiz)}
}

func (r *quizRepo) FindByID(_ context.Context, id uuid.UUID) (*quizdomain.Quiz, error) {
	q, ok := r.items[id]
	if !ok {
		return nil, errors.New("quiz not found")
	}
	return q, nil
}

func (r *quizRepo) Save(_ context.Context, q *quizdomain.Quiz) error {
	r.items[q.ID()] = q
	r.saved = q
	return nil
}

func (r *quizRepo) DeleteByID(_ context.Context, id uuid.UUID) error {
	delete(r.items, id)
	return nil
}

type bankInspector struct {
	counts       map[uuid.UUID]int
	belongsToAll bool
}

func (i *bankInspector) CountQuestionsInBank(_ context.Context, bankID uuid.UUID) (int, error) {
	return i.counts[bankID], nil
}

func (i *bankInspector) QuestionsBelongToBank(_ context.Context, _ uuid.UUID, _ []uuid.UUID) (bool, error) {
	return i.belongsToAll, nil
}

func mustStoredQuiz(t *testing.T, repo *quizRepo, bankID uuid.UUID) *quizdomain.Quiz {
	t.Helper()

	inspector := &bankInspector{counts: map[uuid.UUID]int{bankID: 1}}
	out, err := NewCreateUseCase(repo, inspector).Execute(context.Background(), CreateInput{
		ActorRole: role.NewAdmin(),
		Title:     "Stored quiz",
		Sources: []SourceInput{{
			BankID:        bankID.String(),
			CriteriaType:  criteria.TypeRandom.String(),
			QuestionCount: 1,
		}},
	})
	if err != nil {
		t.Fatalf("create stored quiz: %v", err)
	}

	id, err := uuid.Parse(out.ID)
	if err != nil {
		t.Fatalf("parse stored quiz id: %v", err)
	}

	return repo.items[id]
}
