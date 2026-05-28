package block

import (
	"testing"

	blocktitle "gitflic.ru/lms/backend/internal/domain/course/block/title"
	"github.com/google/uuid"
)

func TestBlockElementIDsWorkflow(t *testing.T) {
	tValue, err := blocktitle.New("Block")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	b, err := New(tValue)
	if err != nil {
		t.Fatalf("new block: %v", err)
	}

	first := uuid.New()
	second := uuid.New()
	third := uuid.New()

	if err := b.AddElementID(first); err != nil {
		t.Fatalf("add first element: %v", err)
	}
	if err := b.AddElementID(second); err != nil {
		t.Fatalf("add second element: %v", err)
	}
	if err := b.AddElementID(third); err != nil {
		t.Fatalf("add third element: %v", err)
	}

	if err := b.MoveElement(2, 0); err != nil {
		t.Fatalf("move element: %v", err)
	}

	ids := b.ElementIDs()
	if len(ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(ids))
	}
	if ids[0] != third || ids[1] != first || ids[2] != second {
		t.Fatalf("unexpected order after move: %v", ids)
	}

	if err := b.RemoveElementID(first); err != nil {
		t.Fatalf("remove element: %v", err)
	}

	ids = b.ElementIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids after remove, got %d", len(ids))
	}
	if ids[0] != third || ids[1] != second {
		t.Fatalf("unexpected order after remove: %v", ids)
	}
}

func TestBlockElementIDsImmutable(t *testing.T) {
	tValue, err := blocktitle.New("Block")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	b, err := New(tValue)
	if err != nil {
		t.Fatalf("new block: %v", err)
	}
	id := uuid.New()
	if err := b.AddElementID(id); err != nil {
		t.Fatalf("add element: %v", err)
	}

	exported := b.ElementIDs()
	exported[0] = uuid.New()

	again := b.ElementIDs()
	if len(again) != 1 || again[0] != id {
		t.Fatalf("element ids must be immutable copy, got %v", again)
	}
}

func TestBlockMoveFromToCompatibility(t *testing.T) {
	tValue, err := blocktitle.New("Block")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	b, err := New(tValue)
	if err != nil {
		t.Fatalf("new block: %v", err)
	}
	first := uuid.New()
	second := uuid.New()
	_ = b.AddElementID(first)
	_ = b.AddElementID(second)

	if err := b.MoveFromTo(1, 0); err != nil {
		t.Fatalf("move from to: %v", err)
	}
	ids := b.ElementIDs()
	if ids[0] != second || ids[1] != first {
		t.Fatalf("unexpected order: %v", ids)
	}
}
