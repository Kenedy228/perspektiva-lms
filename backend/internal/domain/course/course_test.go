package course

import (
	"testing"

	coursetitle "gitflic.ru/lms/backend/internal/domain/course/title"
	"github.com/google/uuid"
)

func TestCourseBlockIDsWorkflow(t *testing.T) {
	tValue, err := coursetitle.New("Course")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	c, err := New(tValue)
	if err != nil {
		t.Fatalf("new course: %v", err)
	}

	first := uuid.New()
	second := uuid.New()
	third := uuid.New()

	if err := c.AddBlockID(first); err != nil {
		t.Fatalf("add first block: %v", err)
	}
	if err := c.AddBlockID(second); err != nil {
		t.Fatalf("add second block: %v", err)
	}
	if err := c.AddBlockID(third); err != nil {
		t.Fatalf("add third block: %v", err)
	}

	if err := c.MoveBlock(2, 0); err != nil {
		t.Fatalf("move block: %v", err)
	}

	ids := c.BlockIDs()
	if len(ids) != 3 {
		t.Fatalf("expected 3 block ids, got %d", len(ids))
	}
	if ids[0] != third || ids[1] != first || ids[2] != second {
		t.Fatalf("unexpected order after move: %v", ids)
	}

	if err := c.RemoveBlockID(first); err != nil {
		t.Fatalf("remove block: %v", err)
	}
	ids = c.BlockIDs()
	if len(ids) != 2 {
		t.Fatalf("expected 2 block ids after remove, got %d", len(ids))
	}
	if ids[0] != third || ids[1] != second {
		t.Fatalf("unexpected order after remove: %v", ids)
	}
}

func TestCourseBlockIDsAreImmutableView(t *testing.T) {
	tValue, err := coursetitle.New("Course")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	c, err := New(tValue)
	if err != nil {
		t.Fatalf("new course: %v", err)
	}
	first := uuid.New()
	second := uuid.New()

	if err := c.AddBlockID(first); err != nil {
		t.Fatalf("add first block id: %v", err)
	}
	if err := c.AddBlockID(second); err != nil {
		t.Fatalf("add second block id: %v", err)
	}

	got := c.BlockIDs()
	got[0] = uuid.New()

	current := c.BlockIDs()
	if len(current) != 2 {
		t.Fatalf("unexpected block ids count: %d", len(current))
	}
	if current[0] != first || current[1] != second {
		t.Fatalf("block ids were mutated through getter: %v", current)
	}
}
