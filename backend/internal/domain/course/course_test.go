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

func TestCourseCompatibilityVersionMethods(t *testing.T) {
	tValue, err := coursetitle.New("Course")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	c, err := New(tValue)
	if err != nil {
		t.Fatalf("new course: %v", err)
	}
	id := uuid.New()

	if err := c.AddVersionID(id); err != nil {
		t.Fatalf("add version id: %v", err)
	}
	if !c.HasVersion(id) {
		t.Fatal("expected version to be present")
	}
	if got := c.VersionIDs(); len(got) != 1 || got[0] != id {
		t.Fatalf("unexpected version ids: %v", got)
	}
	if err := c.RemoveVersionID(id); err != nil {
		t.Fatalf("remove version id: %v", err)
	}
	if c.HasVersion(id) {
		t.Fatal("expected version to be removed")
	}
}
