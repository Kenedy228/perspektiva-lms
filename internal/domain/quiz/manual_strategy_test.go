package quiz

import (
	"errors"
	"slices"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/short"
	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		qSize  int
		addNil bool
		err    error
	}{
		{
			name:   "contains nil question",
			qSize:  10,
			addNil: true,
			err:    ErrNilQuestion,
		},
		{
			name:   "empty questions",
			qSize:  0,
			addNil: false,
			err:    ErrEmptyQuestions,
		},
		{
			name:   "correct",
			qSize:  10,
			addNil: false,
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := tt.qSize
			if tt.addNil {
				count--
			}

			questions, err := generateQuestions(count)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			strategy, err := New(questions)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if count := strategy.CountQuestions(); count != tt.qSize {
				t.Errorf("expected count %v, got %v", tt.qSize, count)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	questions, err := generateQuestions(5)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	}

	strategy, err := New(questions)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	}

	if slices.Compare(questions, strategy.Shuffle()) == 0 {
		t.Errorf("expected init questions and shuffled be different")
	}
}

func generateQuestions(size int) ([]question.Question, error) {
	questions := make([]question.Question, 0, size)

	for range size {
		q, err := short.New(&short.Params{
			Text:           "text",
			Image:          uuid.Nil,
			Answers:        []string{"answer1"},
			AllowDuplicate: false,
		})

		if err != nil {
			return nil, err
		}

		questions = append(questions, q)
	}

	return questions, nil
}
