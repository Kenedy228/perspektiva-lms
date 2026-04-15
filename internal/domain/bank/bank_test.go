package bank

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		title string
		err   error
	}{
		{
			name:  "empty title",
			title: "",
			err:   ErrEmptyTitle,
		},
		{
			name:  "whitespaces title",
			title: " ",
			err:   ErrEmptyTitle,
		},
		{
			name:  "valid title",
			title: "valid",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(tt.title)

			if !errors.Is(tt.err, err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if b != nil {
				if b.ID() == uuid.Nil {
					t.Errorf("expected id non nil, got nil")
				}

				if b.Title() != tt.title {
					t.Errorf("Expected title %v, got %v", tt.title, b.Title())
				}

				if len(b.Questions()) != 0 {
					t.Errorf("expected len of questions 0, got %d", len(b.Questions()))
				}

				if b.CreatedAt() != b.UpdatedAt() {
					t.Errorf("createdAt should equal to updatedAt")
				}
			}
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		err  error
	}{
		{
			name: "empty new",
			old:  "title",
			new:  "",
			err:  ErrEmptyTitle,
		},
		{
			name: "whitespaces new",
			old:  "title",
			new:  "  ",
			err:  ErrEmptyTitle,
		},
		{
			name: "valid new",
			old:  "title",
			new:  "new title",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(tt.old)

			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			oldUpdatedAt := b.UpdatedAt()
			err = b.Rename(tt.new)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if b.Title() == tt.old {
					t.Errorf("expected rename title, but got title unchanged")
				}

				if b.Title() != tt.new {
					t.Errorf("expected rename title, got %v", b.Title())
				}

				if oldUpdatedAt.After(b.UpdatedAt()) {
					t.Errorf("expected updatedAt change")
				}
			}
		})
	}
}

func TestAddQuestionsWithNoErrors(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		questions []uuid.UUID
		total     int
	}{
		{
			name:      "no errors questions",
			title:     "title",
			questions: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			total:     3,
		},
		{
			name:      "no questions",
			title:     "title",
			questions: []uuid.UUID{},
			total:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := New(tt.title)

			total, err := b.AddQuestions(tt.questions...)

			if err != nil {
				t.Errorf("expected no errors, got %v", err)
			}

			if total != tt.total {
				t.Errorf("expected total %v, got %v", tt.total, total)
			}

			if len(b.Questions()) != total {
				t.Errorf("expected len %d, got %d", total, len(b.Questions()))
			}

			if len(b.Questions()) != 0 {
				firstQCopy := b.Questions()
				secondQCopy := b.Questions()

				if &firstQCopy[0] == &secondQCopy[0] {
					t.Errorf("expected slices with different addresses")
				}
			}
		})
	}
}

func TestAddQuestionsWithError(t *testing.T) {
	bank, _ := New("title")
	q := uuid.New()

	total, err := bank.AddQuestions(q, q)

	if !errors.Is(err, ErrQuestionDuplicate) {
		t.Errorf("expected err %v, got %v", ErrQuestionDuplicate, err)
	}

	if total != 0 {
		t.Errorf("expected total 0, got %v", total)
	}
}

func BenchmarkAddQuestionsWithNoError(b *testing.B) {
	sizes := []int{10, 1000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size is %d", size), func(b *testing.B) {
			questions := make([]uuid.UUID, size)

			for i := range questions {
				questions[i] = uuid.New()
			}

			bank, _ := New("title")

			for b.Loop() {
				bank.AddQuestions(questions...)
			}
		})
	}
}

func BenchmarkAddQuestionsWithError(b *testing.B) {
	sizes := []int{9, 999, 99999}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size is %d", size), func(b *testing.B) {
			questions := make([]uuid.UUID, size+1)

			for i := range questions {
				questions[i] = uuid.New()
			}

			questions[len(questions)-1] = questions[len(questions)-2]

			bank, _ := New("title")

			for b.Loop() {
				bank.AddQuestions(questions...)
			}
		})
	}
}

func TestRemoveQuestions(t *testing.T) {
	bank, _ := New("title")
	toDel := uuid.New()

	bank.AddQuestions(uuid.New(), uuid.New(), toDel)
	bank.RemoveQuestions(toDel)

	questions := bank.Questions()

	if slices.Index(questions, toDel) != -1 {
		t.Errorf("expected to delete %v, got present", toDel)
	}
}

func TestClearQuestions(t *testing.T) {
	bank, _ := New("title")

	bank.AddQuestions(uuid.New(), uuid.New(), uuid.New())

	bank.ClearQuestions()

	if length := len(bank.Questions()); length != 0 {
		t.Errorf("expected to clear questions, got len %d", length)
	}
}
