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

func TestCourseAddBlockIDRejectsNil(t *testing.T) {
	c := mustNewCourse(t)
	if err := c.AddBlockID(uuid.Nil); err == nil {
		t.Fatal("expected error for nil block id")
	}
}

func TestCourseAddBlockIDRejectsDuplicate(t *testing.T) {
	c := mustNewCourse(t)
	blockID := uuid.New()
	if err := c.AddBlockID(blockID); err != nil {
		t.Fatalf("add block: %v", err)
	}
	if err := c.AddBlockID(blockID); err == nil {
		t.Fatal("expected error for duplicate block id")
	}
}

func TestCourseAddBlockIDRejectsLimit(t *testing.T) {
	c := mustNewCourse(t)
	for i := 0; i < blocksLimit; i++ {
		if err := c.AddBlockID(uuid.New()); err != nil {
			t.Fatalf("add block %d: %v", i, err)
		}
	}
	if err := c.AddBlockID(uuid.New()); err == nil {
		t.Fatal("expected error when exceeding block limit")
	}
}

func TestCourseRemoveBlockIDRejectsNonExistent(t *testing.T) {
	c := mustNewCourse(t)
	if err := c.RemoveBlockID(uuid.New()); err == nil {
		t.Fatal("expected error for non-existent block id")
	}
}

func TestCourseRemoveBlockIDRejectsNil(t *testing.T) {
	c := mustNewCourse(t)
	if err := c.RemoveBlockID(uuid.Nil); err == nil {
		t.Fatal("expected error for nil block id")
	}
}

func TestCourseMoveBlockOutOfBounds(t *testing.T) {
	c := mustNewCourse(t)
	_ = c.AddBlockID(uuid.New())
	_ = c.AddBlockID(uuid.New())

	if err := c.MoveBlock(0, 5); err == nil {
		t.Fatal("expected error for out of bounds to")
	}
	if err := c.MoveBlock(5, 0); err == nil {
		t.Fatal("expected error for out of bounds from")
	}
	if err := c.MoveBlock(-1, 0); err == nil {
		t.Fatal("expected error for negative from")
	}
}

func TestCourseMoveBlockSamePosition(t *testing.T) {
	c := mustNewCourse(t)
	_ = c.AddBlockID(uuid.New())

	if err := c.MoveBlock(0, 0); err != nil {
		t.Fatalf("move to same position: %v", err)
	}
}

func TestCourseClonePreservesBlockIDs(t *testing.T) {
	c := mustNewCourse(t)
	first := uuid.New()
	_ = c.AddBlockID(first)

	cloned := c.Clone()
	ids := cloned.BlockIDs()
	if len(ids) != 1 || ids[0] != first {
		t.Fatal("cloned course has different block ids")
	}

	ids[0] = uuid.New()
	current := cloned.BlockIDs()
	if current[0] != first {
		t.Fatal("clone block ids were not immutable")
	}
}

func TestCourseRestoreRejectsNilID(t *testing.T) {
	tValue, _ := coursetitle.New("Course")
	_, err := Restore(uuid.Nil, tValue, nil)
	if err == nil {
		t.Fatal("expected error for nil id during restore")
	}
}

func TestCourseRestoreRejectsNilBlockIDs(t *testing.T) {
	tValue, _ := coursetitle.New("Course")
	_, err := Restore(uuid.New(), tValue, []uuid.UUID{uuid.Nil})
	if err == nil {
		t.Fatal("expected error for nil block id in restore")
	}
}

func TestCourseChangeTitleRejectsInvalid(t *testing.T) {
	c := mustNewCourse(t)
	if err := c.ChangeTitle(coursetitle.Title{}); err == nil {
		t.Fatal("expected error for invalid title")
	}
}

func mustNewCourse(t *testing.T) *Course {
	t.Helper()
	tValue, err := coursetitle.New("Course")
	if err != nil {
		t.Fatalf("new title: %v", err)
	}
	c, err := New(tValue)
	if err != nil {
		t.Fatalf("new course: %v", err)
	}
	return c
}

func TestCourseHasNoVersionMethods(t *testing.T) {
	c := mustNewCourse(t)
	_ = c.BlockIDs()
	_ = c.AddBlockID(uuid.New())
	_ = c.RemoveBlockID(c.BlockIDs()[0])
	_ = c.MoveBlock(0, 0)

	var iface interface{} = c
	if _, ok := iface.(interface{ VersionIDs() []uuid.UUID }); ok {
		t.Fatal("Course must not expose VersionIDs method")
	}
}
