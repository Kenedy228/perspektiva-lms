package enrollment

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewAtUsesReferenceDate(t *testing.T) {
	now := time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC)
	_, err := NewAt(uuid.New(), uuid.New(), uuid.New(), now.AddDate(0, 0, -1), now, now)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid activation date, got %v", err)
	}
}

func TestRestoreAllowsHistoricalEnrollment(t *testing.T) {
	activatedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	deactivatedAt := time.Date(2024, 12, 31, 12, 0, 0, 0, time.UTC)

	e, err := Restore(uuid.New(), uuid.New(), uuid.New(), uuid.New(), activatedAt, deactivatedAt)
	if err != nil {
		t.Fatalf("restore enrollment: %v", err)
	}
	if e.Status(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)) != StatusActive {
		t.Fatal("expected active historical enrollment")
	}
}

func TestDeactivate(t *testing.T) {
	now := time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)
	e, err := NewAt(uuid.New(), uuid.New(), uuid.New(), now, now.AddDate(0, 1, 0), now)
	if err != nil {
		t.Fatalf("create enrollment: %v", err)
	}

	err = e.Deactivate(now.AddDate(0, 0, -1))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid deactivation date, got %v", err)
	}

	if err := e.Deactivate(now.AddDate(0, 0, 1)); err != nil {
		t.Fatalf("deactivate enrollment: %v", err)
	}
	if e.DeactivatedAt() != now.AddDate(0, 0, 1) {
		t.Fatalf("unexpected deactivation date: %v", e.DeactivatedAt())
	}
}
