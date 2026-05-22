package progress

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProgressMarkersAndCompletion(t *testing.T) {
	p, err := New(uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	err = p.MarkElement(uuid.New(), MarkerRead, time.Now())
	if err != nil {
		t.Fatalf("mark progress: %v", err)
	}
	if p.CompletedCount() != 1 {
		t.Fatalf("expected one completed item, got %d", p.CompletedCount())
	}
	if p.CompletionPercent(4) != 25 {
		t.Fatalf("expected 25 percent, got %d", p.CompletionPercent(4))
	}
}

func TestProgressRejectsInvalidMarker(t *testing.T) {
	p, err := New(uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	err = p.MarkElement(uuid.New(), MarkerType("unknown"), time.Now())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid marker, got %v", err)
	}
}
