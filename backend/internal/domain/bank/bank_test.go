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

func TestAddQuestions_RejectsNilUUID(t *testing.T) {
	b := mustBank(t)

	err := b.AddQuestions(uuid.Nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("ожидалась ошибка ErrInvalid при добавлении пустого UUID, получено: %v", err)
	}
}

func TestRemoveQuestions_RejectsNilUUID(t *testing.T) {
	b := mustBank(t)

	err := b.RemoveQuestions(uuid.Nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("ожидалась ошибка ErrInvalid при удалении пустого UUID, получено: %v", err)
	}
}

func TestRemoveQuestions_RejectsDuplicatesInRequest(t *testing.T) {
	b := mustBank(t)
	id := uuid.New()
	if err := b.AddQuestions(id); err != nil {
		t.Fatalf("add question: %v", err)
	}

	err := b.RemoveQuestions(id, id)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("ожидалась ошибка ErrInvalid при дублировании в списке удаления, получено: %v", err)
	}
}

func TestAddQuestions_NoOp_WhenEmpty(t *testing.T) {
	b := mustBank(t)

	if err := b.AddQuestions(); err != nil {
		t.Fatalf("добавление пустого списка не должно возвращать ошибку: %v", err)
	}
	if b.CountQuestions() != 0 {
		t.Fatal("количество вопросов должно остаться 0")
	}
}

func TestClearQuestions(t *testing.T) {
	b := mustBank(t)
	if err := b.AddQuestions(uuid.New(), uuid.New()); err != nil {
		t.Fatalf("add questions: %v", err)
	}

	b.ClearQuestions()

	if b.CountQuestions() != 0 {
		t.Fatalf("ожидалось CountQuestions()=0 после ClearQuestions, получено %d", b.CountQuestions())
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
