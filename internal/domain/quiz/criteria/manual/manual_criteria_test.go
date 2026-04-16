package manual

import (
	"errors"
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
)

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
			params := Params{
				QuestionIDs: tt.ids,
			}

			c, err := NewManualCriteria(params)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if c.Type() != criteria.CriteriaTypeManual {
					t.Errorf("expected type %v, got %v", criteria.CriteriaTypeManual, c.Type())
				}

				if c.QuestionCount() != len(tt.ids) {
					t.Errorf("expected len %d, got %d", len(tt.ids), c.QuestionCount())
				}

				if c.QuestionCount() > 0 {
					ids := c.QuestionIDs()
					if &tt.ids[0] == &ids[0] {
						t.Errorf("expected to find copy of slice, got the same slice")
					}
				}
			}
		})
	}
}
