package commands

import (
	"context"
	"errors"
	"testing"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	bankdomain "gitflic.ru/lms/backend/internal/domain/bank"
	"gitflic.ru/lms/backend/internal/domain/bank/title"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestCreateAllowsAdminAndCreator(t *testing.T) {
	ctx := context.Background()

	for _, actor := range []role.Role{role.NewAdmin(), role.NewCreator()} {
		repo := newBankRepo()
		audit := &auditRecorder{}

		out, err := NewCreateUseCase(repo, audit).Execute(ctx, CreateInput{
			ActorRole: actor,
			Title:     "Bank",
		})
		if err != nil {
			t.Fatalf("create bank: %v", err)
		}
		if out.ID == "" || repo.saved == nil {
			t.Fatal("expected bank to be saved")
		}
		if audit.last.Action != bankports.AuditActionCreate {
			t.Fatalf("expected create audit action, got %s", audit.last.Action)
		}
	}
}

func TestCreateRejectsStudent(t *testing.T) {
	repo := newBankRepo()

	_, err := NewCreateUseCase(repo, nil).Execute(context.Background(), CreateInput{
		ActorRole: role.NewStudent(),
		Title:     "Bank",
	})
	if !errors.Is(err, common.ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
	if repo.saved != nil {
		t.Fatal("student-created bank must not be saved")
	}
}

func TestAddAndRemoveQuestionsUseCases(t *testing.T) {
	ctx := context.Background()
	repo := newBankRepo()
	b := mustStoredBank(t, repo)
	first := uuid.New()
	second := uuid.New()

	err := NewAddQuestionsUseCase(repo, nil).Execute(ctx, QuestionIDsInput{
		ActorRole:   role.NewCreator(),
		BankID:      b.ID().String(),
		QuestionIDs: []string{first.String(), second.String()},
	})
	if err != nil {
		t.Fatalf("add questions: %v", err)
	}
	if repo.saved.CountQuestions() != 2 {
		t.Fatalf("expected 2 questions, got %d", repo.saved.CountQuestions())
	}

	err = NewRemoveQuestionsUseCase(repo, nil).Execute(ctx, QuestionIDsInput{
		ActorRole:   role.NewAdmin(),
		BankID:      b.ID().String(),
		QuestionIDs: []string{first.String()},
	})
	if err != nil {
		t.Fatalf("remove questions: %v", err)
	}
	if repo.saved.HasQuestion(first) {
		t.Fatal("expected question to be removed")
	}
}

func TestRenameClearAndDeleteUseCases(t *testing.T) {
	ctx := context.Background()
	repo := newBankRepo()
	b := mustStoredBank(t, repo)
	questionID := uuid.New()
	if err := b.AddQuestions(questionID); err != nil {
		t.Fatalf("add question: %v", err)
	}

	if err := NewRenameUseCase(repo, nil).Execute(ctx, RenameInput{
		ActorRole: role.NewCreator(),
		BankID:    b.ID().String(),
		Title:     "Renamed",
	}); err != nil {
		t.Fatalf("rename bank: %v", err)
	}
	if repo.saved.Title().Value() != "Renamed" {
		t.Fatalf("expected renamed bank, got %q", repo.saved.Title().Value())
	}

	if err := NewClearQuestionsUseCase(repo, nil).Execute(ctx, BankIDInput{
		ActorRole: role.NewAdmin(),
		BankID:    b.ID().String(),
	}); err != nil {
		t.Fatalf("clear questions: %v", err)
	}
	if repo.saved.CountQuestions() != 0 {
		t.Fatalf("expected cleared questions, got %d", repo.saved.CountQuestions())
	}

	if err := NewDeleteUseCase(repo, nil).Execute(ctx, BankIDInput{
		ActorRole: role.NewAdmin(),
		BankID:    b.ID().String(),
	}); err != nil {
		t.Fatalf("delete bank: %v", err)
	}
	if _, ok := repo.items[b.ID()]; ok {
		t.Fatal("expected bank to be deleted")
	}
}

type bankRepo struct {
	items map[uuid.UUID]*bankdomain.Bank
	saved *bankdomain.Bank
}

func newBankRepo() *bankRepo {
	return &bankRepo{items: make(map[uuid.UUID]*bankdomain.Bank)}
}

func (r *bankRepo) FindByID(_ context.Context, id uuid.UUID) (*bankdomain.Bank, error) {
	b, ok := r.items[id]
	if !ok {
		return nil, errors.New("bank not found")
	}
	return b, nil
}

func (r *bankRepo) Save(_ context.Context, b *bankdomain.Bank) error {
	r.items[b.ID()] = b
	r.saved = b
	return nil
}

func (r *bankRepo) DeleteByID(_ context.Context, id uuid.UUID) error {
	delete(r.items, id)
	return nil
}

type auditRecorder struct {
	last bankports.AuditEvent
}

func (r *auditRecorder) RecordBankAudit(_ context.Context, event bankports.AuditEvent) error {
	r.last = event
	return nil
}

func mustStoredBank(t *testing.T, repo *bankRepo) *bankdomain.Bank {
	t.Helper()

	tl, err := title.New("Bank")
	if err != nil {
		t.Fatalf("create title: %v", err)
	}
	b, err := bankdomain.New(tl)
	if err != nil {
		t.Fatalf("create bank: %v", err)
	}
	if err := repo.Save(context.Background(), b); err != nil {
		t.Fatalf("save bank: %v", err)
	}
	return b
}
