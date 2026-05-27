package item

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	qtitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
)

type mockQuestion struct {
	id    uuid.UUID
	title qtitle.Title
}

func (q *mockQuestion) ID() uuid.UUID                    { return q.id }
func (q *mockQuestion) Title() qtitle.Title              { return q.title }
func (q *mockQuestion) Instruction() string              { return "инструкция" }
func (q *mockQuestion) Type() question.Type              { return question.TypeShort }
func (q *mockQuestion) ChangeTitle(t qtitle.Title) error { q.title = t; return nil }
func (q *mockQuestion) Clone() question.Question {
	copy := *q
	return &copy
}

func TestNew(t *testing.T) {
	ttl, _ := qtitle.New("Вопрос")
	q := &mockQuestion{id: uuid.New(), title: ttl}

	tests := []struct {
		name    string
		q       question.Question
		wantErr error
	}{
		{name: "ok", q: q},
		{name: "nil question", q: nil, wantErr: ErrInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it, err := New(tt.q)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if it.ID() != q.id {
				t.Fatalf("unexpected id: %s", it.ID())
			}
		})
	}
}

func TestSnapshotReturnsClone(t *testing.T) {
	ttl, _ := qtitle.New("Вопрос")
	q := &mockQuestion{id: uuid.New(), title: ttl}

	it, err := New(q)
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	snap := it.Snapshot()
	newTitle, _ := qtitle.New("Новое имя")
	if err := snap.ChangeTitle(newTitle); err != nil {
		t.Fatalf("change snapshot title: %v", err)
	}

	again := it.Snapshot()
	if again.Title().Value() == newTitle.Value() {
		t.Fatal("snapshot leaked mutable reference")
	}
}
