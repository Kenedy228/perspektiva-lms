package answer

import (
	"errors"
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

type mockAnswer struct {
	value string
}

func (a *mockAnswer) IsEmpty() bool { return a.value == "" }
func (a *mockAnswer) Clone() question.Answer {
	copy := *a
	return &copy
}

func TestNew(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		questionID uuid.UUID
		ans        question.Answer
		answeredAt time.Time
		wantErr    error
	}{
		{name: "ok", questionID: uuid.New(), ans: &mockAnswer{value: "x"}, answeredAt: now},
		{name: "zero question id", questionID: uuid.Nil, ans: &mockAnswer{value: "x"}, answeredAt: now, wantErr: ErrInvalid},
		{name: "nil answer", questionID: uuid.New(), ans: nil, answeredAt: now, wantErr: ErrInvalid},
		{name: "zero answeredAt", questionID: uuid.New(), ans: &mockAnswer{value: "x"}, answeredAt: time.Time{}, wantErr: ErrInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := New(tt.questionID, tt.ans, tt.answeredAt)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if entry.QuestionID() != tt.questionID {
				t.Fatalf("unexpected question id: %s", entry.QuestionID())
			}
		})
	}
}

func TestAnswerReturnsClone(t *testing.T) {
	id := uuid.New()
	ans := &mockAnswer{value: "исходный"}
	entry, err := New(id, ans, time.Now())
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}

	fromEntry := entry.Answer().(*mockAnswer)
	fromEntry.value = "изменен"

	again := entry.Answer().(*mockAnswer)
	if again.value == "изменен" {
		t.Fatal("entry leaked mutable answer")
	}
}
