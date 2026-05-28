package attempt

import (
	"math/rand"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	questiontitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
)

func TestMathRandQuestionShufflerShufflesCopy(t *testing.T) {
	first := shufflerQuestion{id: uuid.New()}
	second := shufflerQuestion{id: uuid.New()}
	third := shufflerQuestion{id: uuid.New()}
	fourth := shufflerQuestion{id: uuid.New()}
	questions := []question.Question{first, second, third, fourth}

	shuffler := NewMathRandQuestionShufflerWithSource(rand.NewSource(1))
	shuffled := shuffler.ShuffleQuestions(questions)

	if len(shuffled) != len(questions) {
		t.Fatalf("expected %d questions, got %d", len(questions), len(shuffled))
	}
	if questions[0].ID() != first.ID() || questions[1].ID() != second.ID() {
		t.Fatal("expected source slice order to stay unchanged")
	}
	if !sameQuestionSet(questions, shuffled) {
		t.Fatal("expected shuffled slice to contain the same questions")
	}
}

func TestNewMathRandQuestionShufflerWithSourceRequiresSource(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil source")
		}
	}()

	_ = NewMathRandQuestionShufflerWithSource(nil)
}

type shufflerQuestion struct {
	id uuid.UUID
}

func (q shufflerQuestion) ID() uuid.UUID { return q.id }

func (q shufflerQuestion) Title() questiontitle.Title {
	t, _ := questiontitle.New("Question")
	return t
}

func (q shufflerQuestion) Instruction() string { return question.TypeShort.DefaultInstruction() }
func (q shufflerQuestion) Type() question.Type { return question.TypeShort }
func (q shufflerQuestion) Clone() question.Question {
	return shufflerQuestion{id: q.id}
}
func (q shufflerQuestion) ChangeTitle(questiontitle.Title) error { return nil }

func sameQuestionSet(left, right []question.Question) bool {
	if len(left) != len(right) {
		return false
	}

	seen := make(map[uuid.UUID]int, len(left))
	for i := range left {
		seen[left[i].ID()]++
	}
	for i := range right {
		seen[right[i].ID()]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}
