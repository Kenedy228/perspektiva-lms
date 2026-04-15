package quiz

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewSource(t *testing.T) {
	tests := []struct {
		name   string
		bankID uuid.UUID
		err    error
	}{
		{
			name:   "nil bank",
			bankID: uuid.Nil,
			err:    ErrNilBank,
		},
		{
			name:   "valid bank",
			bankID: uuid.New(),
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewRandomCriteria(10)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			s, err := NewSource(tt.bankID, c)
			if !errors.Is(tt.err, err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if s.BankID() != tt.bankID {
					t.Errorf("expected bank id %v, got %v", tt.bankID, s.BankID())
				}

				if s.Criteria().Type() != c.Type() {
					t.Errorf("expected criteria type %v, got %v", c.Type(), s.Criteria().Type())
				}
			}
		})
	}
}

func TestNewSourceWithNilCriteria(t *testing.T) {
	_, err := NewSource(uuid.New(), nil)
	if err != ErrNilCriteria {
		t.Errorf("expected err %v, got %v", ErrNilCriteria, err)
	}
}
