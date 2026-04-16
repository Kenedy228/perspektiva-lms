package dob

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("should not return err", func(t *testing.T) {
		now := time.Now()
		dob, err := NewDateOfBirth(now)
		if err != nil {
			t.Errorf("expected err nil, got %v", err)
		}

		if err == nil {
			if dob.Date() != now {
				t.Errorf("expected Date %v, got %v", now, dob.Date())
			}
		}
	})

	t.Run("should throw err", func(t *testing.T) {
		now := time.Now()
		now = now.Add(time.Hour)

		_, err := NewDateOfBirth(now)

		if !errors.Is(err, ErrInvalidDob) {
			t.Errorf("expected err %v, got %v", ErrInvalidDob, err)
		}
	})
}

func TestAgeAt(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		dob := time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)

		d, err := NewDateOfBirth(dob)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		age, err := d.AgeAt(time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC))
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		if age != 26 {
			t.Errorf("expected age %v, got %v", 26, age)
		}
	})

	t.Run("invalid date", func(t *testing.T) {
		now := time.Now()

		d, err := NewDateOfBirth(now)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		at := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		_, err = d.AgeAt(at)
		if !errors.Is(err, ErrInvalidAt) {
			t.Errorf("expected err %v, got %v", ErrInvalidAt, err)
		}
	})
}
