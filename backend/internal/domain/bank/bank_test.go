package bank

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/bank/title"
	"github.com/google/uuid"
)

func TestNewRequiresTitle(t *testing.T) {
	_, err := New(title.Title{})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid bank error, got %v", err)
	}
}

func TestRestoreValidatesIDTitleAndQuestions(t *testing.T) {
	tl := mustTitle(t)

	_, err := Restore(uuid.Nil, tl, nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid bank id error, got %v", err)
	}

	questionID := uuid.New()
	_, err = Restore(uuid.New(), tl, []uuid.UUID{questionID, questionID})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected duplicate question error, got %v", err)
	}

	b, err := Restore(uuid.New(), tl, []uuid.UUID{questionID})
	if err != nil {
		t.Fatalf("restore bank: %v", err)
	}
	if !b.HasQuestion(questionID) {
		t.Fatal("expected restored question")
	}
}

func TestRenameRequiresTitle(t *testing.T) {
	b := mustBank(t)

	err := b.Rename(title.Title{})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid title error, got %v", err)
	}
}

func TestAddAndRemoveQuestions(t *testing.T) {
	b := mustBank(t)
	first := uuid.New()
	second := uuid.New()

	if err := b.AddQuestions(first, second); err != nil {
		t.Fatalf("add questions: %v", err)
	}
	if b.CountQuestions() != 2 {
		t.Fatalf("expected 2 questions, got %d", b.CountQuestions())
	}

	err := b.AddQuestions(first)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected duplicate add error, got %v", err)
	}

	err = b.RemoveQuestions(uuid.New())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected missing remove error, got %v", err)
	}

	if err := b.RemoveQuestions(first); err != nil {
		t.Fatalf("remove question: %v", err)
	}
	if b.HasQuestion(first) {
		t.Fatal("expected question to be removed")
	}
}

func mustBank(t *testing.T) *Bank {
	t.Helper()

	b, err := New(mustTitle(t))
	if err != nil {
		t.Fatalf("create bank: %v", err)
	}
	return b
}

func mustTitle(t *testing.T) title.Title {
	t.Helper()

	tl, err := title.New("Question bank")
	if err != nil {
		t.Fatalf("create title: %v", err)
	}
	return tl
}
