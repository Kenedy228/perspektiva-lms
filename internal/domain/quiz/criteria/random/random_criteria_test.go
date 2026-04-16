package random

import (
	"errors"
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
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
			params := Params{
				Count: tt.count,
			}
			c, err := NewRandomCriteria(params)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if c.Type() != criteria.CriteriaTypeRandom {
					t.Errorf("expected type %v, got %v", criteria.CriteriaTypeRandom, c.Type())
				}

				if c.QuestionCount() != tt.count {
					t.Errorf("expected count %v, got %v", tt.count, c.QuestionCount())
				}
			}
		})
	}
}
