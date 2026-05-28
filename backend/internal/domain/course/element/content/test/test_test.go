package test

import (
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		quizID  uuid.UUID
		wantErr bool
	}{
		{name: "valid quiz id", quizID: uuid.New(), wantErr: false},
		{name: "nil quiz id", quizID: uuid.Nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.quizID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	c, err := New(uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if ct := c.ContentType(); ct.String() != "test" {
		t.Fatalf("unexpected content type: %q", ct)
	}
}

func TestIsInteractive(t *testing.T) {
	c, err := New(uuid.New())
	if err != nil {
		t.Fatal(err)
	}
	if !c.IsInteractive() {
		t.Fatal("expected test content to be interactive")
	}
}

func TestClone(t *testing.T) {
	id := uuid.New()
	c, err := New(id)
	if err != nil {
		t.Fatal(err)
	}
	cloned := c.Clone()
	switch ct := cloned.(type) {
	case Content:
		if ct.QuizID() != id {
			t.Fatal("cloned content has different quiz id")
		}
	default:
		t.Fatalf("unexpected cloned type: %T", cloned)
	}
}
