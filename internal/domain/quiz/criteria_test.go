package quiz

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestRandomCriteria(t *testing.T) {
	tests := []struct {
		name  string
		count int
		err   error
	}{
		{
			name:  "negative count",
			count: -1,
			err:   ErrInvalidCount,
		},
		{
			name:  "zero count",
			count: 0,
			err:   ErrInvalidCount,
		},
		{
			name:  "valid count",
			count: 1,
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewRandomCriteria(tt.count)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if c.Type() != CriteriaTypeRandom {
					t.Errorf("expected type %v, got %v", CriteriaTypeRandom, c.Type())
				}

				if c.Count() != tt.count {
					t.Errorf("expected count %v, got %v", tt.count, c.Count())
				}
			}
		})
	}
}

func TestManualCriteria(t *testing.T) {
	tests := []struct {
		name string
		ids  []uuid.UUID
		err  error
	}{
		{
			name: "empty ids",
			ids:  []uuid.UUID{},
			err:  ErrEmptyQuestions,
		},
		{
			name: "nil id",
			ids:  []uuid.UUID{uuid.Nil},
			err:  ErrNilQuestion,
		},
		{
			name: "valid ids",
			ids:  []uuid.UUID{uuid.New(), uuid.New()},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewManualCriteria(tt.ids)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if c.Type() != CriteriaTypeManual {
					t.Errorf("expected type %v, got %v", CriteriaTypeManual, c.Type())
				}

				ids := c.QuestionIDs()

				if len(ids) != len(tt.ids) {
					t.Errorf("expected len %d, got %d", len(tt.ids), len(ids))
				}

				if len(ids) > 0 {
					if &tt.ids[0] == &ids[0] {
						t.Errorf("expected to find copy of slice, got the same slice")
					}
				}
			}
		})
	}
}
