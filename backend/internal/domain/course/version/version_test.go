package version

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/course/version/title"
	"github.com/google/uuid"
)

func TestPublishRequiresBlocksAndLocksEditing(t *testing.T) {
	v, err := New(mustTitle(t))
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	err = v.Publish()
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid publish without blocks, got %v", err)
	}

	if err := v.AddBlockID(uuid.New()); err != nil {
		t.Fatalf("add block: %v", err)
	}
	if err := v.Publish(); err != nil {
		t.Fatalf("publish version: %v", err)
	}

	err = v.AddBlockID(uuid.New())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected editing published version to fail, got %v", err)
	}
}

func TestRestoreValidatesStatusAndBlocks(t *testing.T) {
	blockID := uuid.New()
	_, err := Restore(uuid.New(), mustTitle(t), Status("unknown"), []uuid.UUID{blockID})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid status, got %v", err)
	}

	_, err = Restore(uuid.New(), mustTitle(t), StatusDraft, []uuid.UUID{blockID, blockID})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected duplicate blocks error, got %v", err)
	}
}

func mustTitle(t *testing.T) title.Title {
	t.Helper()
	tl, err := title.New("Version")
	if err != nil {
		t.Fatalf("create title: %v", err)
	}
	return tl
}
