package progress

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProgressMarkersAndCompletion(t *testing.T) {
	p, err := New(uuid.New())
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
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	err = p.MarkElement(uuid.New(), MarkerType("unknown"), time.Now())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid marker, got %v", err)
	}
}

func TestProgressMarkUnmarkAndCheckCompleted(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	elementID := uuid.New()
	if err := p.MarkCompleted(elementID); err != nil {
		t.Fatalf("mark completed: %v", err)
	}
	if !p.IsCompleted(elementID) {
		t.Fatal("expected element to be completed")
	}
	if err := p.UnmarkCompleted(elementID); err != nil {
		t.Fatalf("unmark completed: %v", err)
	}
	if p.IsCompleted(elementID) {
		t.Fatal("expected element to be not completed")
	}
}

func TestProgressDuplicateMark(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	elementID := uuid.New()
	if err := p.MarkCompleted(elementID); err != nil {
		t.Fatalf("first mark: %v", err)
	}
	if err := p.MarkCompleted(elementID); err != nil {
		t.Fatalf("duplicate mark: %v", err)
	}
	if p.CompletedCount() != 1 {
		t.Fatalf("duplicate mark should not increase count, got %d", p.CompletedCount())
	}
}

func TestProgressPercentEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		markersCount int
		totalTracked int
		expectedPct  int
	}{
		{name: "zero tracked items", markersCount: 0, totalTracked: 0, expectedPct: 0},
		{name: "negative tracked items", markersCount: 5, totalTracked: -1, expectedPct: 0},
		{name: "full completion", markersCount: 5, totalTracked: 5, expectedPct: 100},
		{name: "over completion cap at 100", markersCount: 7, totalTracked: 5, expectedPct: 100},
		{name: "50 percent", markersCount: 2, totalTracked: 4, expectedPct: 50},
		{name: "zero of N", markersCount: 0, totalTracked: 3, expectedPct: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(uuid.New())
			if err != nil {
				t.Fatalf("create progress: %v", err)
			}
			for i := 0; i < tt.markersCount; i++ {
				if err := p.MarkCompleted(uuid.New()); err != nil {
					t.Fatalf("mark element: %v", err)
				}
			}
			if pct := p.CompletionPercent(tt.totalTracked); pct != tt.expectedPct {
				t.Fatalf("expected %d%%, got %d%%", tt.expectedPct, pct)
			}
		})
	}
}

func TestProgressNoVersionID(t *testing.T) {
	enrollmentID := uuid.New()
	p, err := New(enrollmentID)
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}
	if p.EnrollmentID() != enrollmentID {
		t.Fatal("enrollment id mismatch")
	}

	r, err := Restore(uuid.New(), enrollmentID, map[uuid.UUID]Marker{})
	if err != nil {
		t.Fatalf("restore progress: %v", err)
	}
	if r.EnrollmentID() != enrollmentID {
		t.Fatal("restored enrollment id mismatch")
	}
}

func TestProgressMarkCompletedWithTime(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}

	elementID := uuid.New()
	at := time.Date(2026, 5, 28, 12, 0, 0, 0, time.UTC)
	if err := p.MarkElement(elementID, MarkerWatched, at); err != nil {
		t.Fatalf("mark watched: %v", err)
	}
	if !p.IsCompleted(elementID) {
		t.Fatal("expected element to be completed")
	}
	if p.CompletedCount() != 1 {
		t.Fatalf("expected 1 completed, got %d", p.CompletedCount())
	}
}

func TestProgressNewRejectsNilEnrollmentID(t *testing.T) {
	_, err := New(uuid.Nil)
	if err == nil {
		t.Fatal("expected error for nil enrollment id")
	}
}

func TestProgressRestoreRejectsNilEnrollmentID(t *testing.T) {
	_, err := Restore(uuid.New(), uuid.Nil, map[uuid.UUID]Marker{})
	if err == nil {
		t.Fatal("expected error for nil enrollment id during restore")
	}
}

func TestProgressMarkElementRejectsNilElementID(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}
	if err := p.MarkElement(uuid.Nil, MarkerRead, time.Now()); err == nil {
		t.Fatal("expected error for nil element id")
	}
}

func TestProgressUnmarkCompletedRejectsNilElementID(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}
	if err := p.UnmarkCompleted(uuid.Nil); err == nil {
		t.Fatal("expected error for nil element id")
	}
}

func TestProgressMarkersAreImmutable(t *testing.T) {
	p, err := New(uuid.New())
	if err != nil {
		t.Fatalf("create progress: %v", err)
	}
	elementID := uuid.New()
	if err := p.MarkCompleted(elementID); err != nil {
		t.Fatalf("mark completed: %v", err)
	}

	markers := p.Markers()
	markers[uuid.New()] = Marker{mType: MarkerQuiz, completedAt: time.Now()}

	currentMarkers := p.Markers()
	if len(currentMarkers) != 1 {
		t.Fatalf("markers were mutated through getter, got %d", len(currentMarkers))
	}
}
