package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	blocktitle "gitflic.ru/lms/backend/internal/domain/course/block/title"
	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/domain/course/title"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestCreateCourseRejectsStudent(t *testing.T) {
	_, err := NewCreateCourseUseCase(newCourseRepo()).Execute(context.Background(), CreateCourseInput{
		ActorRole: role.NewStudent(),
		Title:     "Course",
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestCreateCourseSuccess(t *testing.T) {
	courses := newCourseRepo()
	out, err := NewCreateCourseUseCase(courses).Execute(context.Background(), CreateCourseInput{
		ActorRole: role.NewAdmin(),
		Title:     "Course",
	})
	if err != nil {
		t.Fatalf("create course: %v", err)
	}
	if out.ID == "" {
		t.Fatal("expected course id")
	}
}

func TestRenameCourseSuccess(t *testing.T) {
	courses := newCourseRepo()
	tValue, _ := title.New("Old")
	c, _ := coursedomain.New(tValue)
	_ = courses.Save(context.Background(), c)

	err := NewRenameCourseUseCase(courses).Execute(context.Background(), RenameCourseInput{
		ActorRole: role.NewCreator(),
		CourseID:  c.ID().String(),
		Title:     "New",
	})
	if err != nil {
		t.Fatalf("rename course: %v", err)
	}
	updated, _ := courses.FindByID(context.Background(), c.ID())
	if updated.Title().Value() != "New" {
		t.Fatalf("expected title 'New', got %q", updated.Title().Value())
	}
}

func TestRenameCourseRejectsStudent(t *testing.T) {
	err := NewRenameCourseUseCase(newCourseRepo()).Execute(context.Background(), RenameCourseInput{
		ActorRole: role.NewStudent(),
		CourseID:  uuid.New().String(),
		Title:     "New",
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestRenameCourseRejectsInvalidID(t *testing.T) {
	err := NewRenameCourseUseCase(newCourseRepo()).Execute(context.Background(), RenameCourseInput{
		ActorRole: role.NewAdmin(),
		CourseID:  "not-a-uuid",
		Title:     "New",
	})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestAddBlockToCourseRejectsStudent(t *testing.T) {
	_, err := NewAddBlockToCourseUseCase(newCourseRepo(), newBlockRepo()).Execute(context.Background(), AddBlockToCourseInput{
		ActorRole: role.NewStudent(),
		CourseID:  uuid.New().String(),
		Title:     "Block",
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestAddBlockToCourseSuccess(t *testing.T) {
	ctx := context.Background()
	courses := newCourseRepo()
	blocks := newBlockRepo()
	tValue, _ := title.New("Course")
	c, _ := coursedomain.New(tValue)
	_ = courses.Save(ctx, c)

	out, err := NewAddBlockToCourseUseCase(courses, blocks).Execute(ctx, AddBlockToCourseInput{
		ActorRole: role.NewCreator(),
		CourseID:  c.ID().String(),
		Title:     "Block",
	})
	if err != nil {
		t.Fatalf("add block: %v", err)
	}
	if out.ID == "" {
		t.Fatal("expected block id")
	}
	updated, _ := courses.FindByID(ctx, c.ID())
	if len(updated.BlockIDs()) != 1 {
		t.Fatalf("expected 1 block, got %d", len(updated.BlockIDs()))
	}
}

func TestMoveCourseBlockRejectsStudent(t *testing.T) {
	err := NewMoveCourseBlockUseCase(newCourseRepo()).Execute(context.Background(), MoveCourseBlockInput{
		ActorRole: role.NewStudent(),
		CourseID:  uuid.New().String(),
		From:      0,
		To:        1,
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestMoveCourseBlockSuccess(t *testing.T) {
	ctx := context.Background()
	courses := newCourseRepo()
	tValue, _ := title.New("Course")
	c, _ := coursedomain.New(tValue)
	_ = c.AddBlockID(uuid.New())
	_ = c.AddBlockID(uuid.New())
	_ = courses.Save(ctx, c)

	originalFirst := c.BlockIDs()[0]
	originalSecond := c.BlockIDs()[1]

	err := NewMoveCourseBlockUseCase(courses).Execute(ctx, MoveCourseBlockInput{
		ActorRole: role.NewAdmin(),
		CourseID:  c.ID().String(),
		From:      1,
		To:        0,
	})
	if err != nil {
		t.Fatalf("move block: %v", err)
	}
	updated, _ := courses.FindByID(ctx, c.ID())
	if updated.BlockIDs()[0] != originalSecond || updated.BlockIDs()[1] != originalFirst {
		t.Fatal("block order not updated")
	}
}

func TestAddElementToBlockRejectsStudent(t *testing.T) {
	_, err := NewAddElementToBlockUseCase(newBlockRepo(), newElementRepo()).Execute(context.Background(), AddElementToBlockInput{
		ActorRole: role.NewStudent(),
		BlockID:   uuid.New().String(),
		Title:     "Element",
		Content: ElementContentInput{
			Type:      "download_file",
			FileName:  "notes.docx",
			SizeBytes: 128,
		},
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestAddElementToBlockSuccess(t *testing.T) {
	ctx := context.Background()
	blocks := newBlockRepo()
	elements := newElementRepo()
	b, _ := block.New(mustBlockTitle(t, "Block"))
	_ = blocks.Save(ctx, b)

	out, err := NewAddElementToBlockUseCase(blocks, elements).Execute(ctx, AddElementToBlockInput{
		ActorRole: role.NewAdmin(),
		BlockID:   b.ID().String(),
		Title:     "Element",
		Content: ElementContentInput{
			Type:   "test",
			QuizID: uuid.New().String(),
		},
	})
	if err != nil {
		t.Fatalf("add element: %v", err)
	}
	if out.ID == "" {
		t.Fatal("expected element id")
	}
	updated, _ := blocks.FindByID(ctx, b.ID())
	if len(updated.ElementIDs()) != 1 {
		t.Fatalf("expected 1 element, got %d", len(updated.ElementIDs()))
	}
}

func TestAddElementToBlockWithCompletionMode(t *testing.T) {
	ctx := context.Background()
	blocks := newBlockRepo()
	elements := newElementRepo()
	b, _ := block.New(mustBlockTitle(t, "Block"))
	_ = blocks.Save(ctx, b)

	_, err := NewAddElementToBlockUseCase(blocks, elements).Execute(ctx, AddElementToBlockInput{
		ActorRole: role.NewAdmin(),
		BlockID:   b.ID().String(),
		Title:     "Element",
		Content: ElementContentInput{
			Type:           "lecture_material",
			FileName:       "lesson.mp4",
			SizeBytes:      1024,
			CompletionMode: "manual",
		},
	})
	if err != nil {
		t.Fatalf("add element with completion mode: %v", err)
	}
}

func TestAddElementToBlockRejectsInvalidContentType(t *testing.T) {
	ctx := context.Background()
	blocks := newBlockRepo()
	elements := newElementRepo()
	b, _ := block.New(mustBlockTitle(t, "Block"))
	_ = blocks.Save(ctx, b)

	_, err := NewAddElementToBlockUseCase(blocks, elements).Execute(ctx, AddElementToBlockInput{
		ActorRole: role.NewAdmin(),
		BlockID:   b.ID().String(),
		Title:     "Element",
		Content: ElementContentInput{
			Type: "unknown_type",
		},
	})
	if err == nil {
		t.Fatal("expected error for unknown content type")
	}
}

func TestMoveBlockElementRejectsStudent(t *testing.T) {
	err := NewMoveBlockElementUseCase(newBlockRepo()).Execute(context.Background(), MoveBlockElementInput{
		ActorRole: role.NewStudent(),
		BlockID:   uuid.New().String(),
		From:      0,
		To:        1,
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestMoveBlockElementSuccess(t *testing.T) {
	ctx := context.Background()
	blocks := newBlockRepo()
	b, _ := block.New(mustBlockTitle(t, "Block"))
	_ = b.AddElementID(uuid.New())
	_ = b.AddElementID(uuid.New())
	_ = blocks.Save(ctx, b)

	originalFirst := b.ElementIDs()[0]
	originalSecond := b.ElementIDs()[1]

	err := NewMoveBlockElementUseCase(blocks).Execute(ctx, MoveBlockElementInput{
		ActorRole: role.NewAdmin(),
		BlockID:   b.ID().String(),
		From:      1,
		To:        0,
	})
	if err != nil {
		t.Fatalf("move element: %v", err)
	}
	updated, _ := blocks.FindByID(ctx, b.ID())
	if updated.ElementIDs()[0] != originalSecond || updated.ElementIDs()[1] != originalFirst {
		t.Fatal("element order not updated")
	}
}

func TestMarkProgressRejectsNonStudent(t *testing.T) {
	err := NewMarkProgressUseCase(newProgressRepo()).Execute(context.Background(), MarkProgressInput{
		ActorRole:    role.NewAdmin(),
		EnrollmentID: uuid.New().String(),
		ElementID:    uuid.New().String(),
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestCourseBlockElementWorkflow(t *testing.T) {
	ctx := context.Background()
	courses := newCourseRepo()
	blocks := newBlockRepo()
	elements := newElementRepo()

	courseTitle, _ := title.New("Course")
	c, _ := coursedomain.New(courseTitle)
	_ = courses.Save(ctx, c)

	blockOut, err := NewAddBlockToCourseUseCase(courses, blocks).Execute(ctx, AddBlockToCourseInput{
		ActorRole: role.NewCreator(),
		CourseID:  c.ID().String(),
		Title:     "Block",
	})
	if err != nil {
		t.Fatalf("add block to course: %v", err)
	}

	if err := NewMoveCourseBlockUseCase(courses).Execute(ctx, MoveCourseBlockInput{
		ActorRole: role.NewAdmin(),
		CourseID:  c.ID().String(),
		From:      0,
		To:        0,
	}); err != nil {
		t.Fatalf("move course block: %v", err)
	}

	elementOut, err := NewAddElementToBlockUseCase(blocks, elements).Execute(ctx, AddElementToBlockInput{
		ActorRole: role.NewAdmin(),
		BlockID:   blockOut.ID,
		Title:     "Material",
		Content: ElementContentInput{
			Type:      "download_file",
			FileName:  "notes.docx",
			SizeBytes: 128,
		},
	})
	if err != nil {
		t.Fatalf("add element to block: %v", err)
	}
	if elementOut.ID == "" {
		t.Fatal("expected element id")
	}

	if err := NewMoveBlockElementUseCase(blocks).Execute(ctx, MoveBlockElementInput{
		ActorRole: role.NewAdmin(),
		BlockID:   blockOut.ID,
		From:      0,
		To:        0,
	}); err != nil {
		t.Fatalf("move block element: %v", err)
	}
}

func TestMarkProgressUseCase(t *testing.T) {
	ctx := context.Background()
	repo := newProgressRepo()
	p, _ := progress.New(uuid.New())
	repo.item = p

	err := NewMarkProgressUseCase(repo).Execute(ctx, MarkProgressInput{
		ActorRole:    role.NewStudent(),
		EnrollmentID: p.EnrollmentID().String(),
		ElementID:    uuid.New().String(),
		MarkerType:   progress.MarkerRead,
		At:           time.Now(),
	})
	if err != nil {
		t.Fatalf("mark progress: %v", err)
	}
	if repo.saved.CompletedCount() != 1 {
		t.Fatalf("expected one marker, got %d", repo.saved.CompletedCount())
	}
}

type courseRepo struct {
	items map[uuid.UUID]*coursedomain.Course
}

func newCourseRepo() *courseRepo { return &courseRepo{items: make(map[uuid.UUID]*coursedomain.Course)} }
func (r *courseRepo) FindByID(_ context.Context, id uuid.UUID) (*coursedomain.Course, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errors.New("course not found")
	}
	return item, nil
}
func (r *courseRepo) Save(_ context.Context, c *coursedomain.Course) error {
	r.items[c.ID()] = c
	return nil
}

type blockRepo struct{ items map[uuid.UUID]*block.Block }

func newBlockRepo() *blockRepo { return &blockRepo{items: make(map[uuid.UUID]*block.Block)} }
func (r *blockRepo) FindByID(_ context.Context, id uuid.UUID) (*block.Block, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errors.New("block not found")
	}
	return item, nil
}
func (r *blockRepo) Save(_ context.Context, b *block.Block) error {
	r.items[b.ID()] = b
	return nil
}

type elementRepo struct {
	items map[uuid.UUID]*elementdomain.Element
}

func newElementRepo() *elementRepo {
	return &elementRepo{items: make(map[uuid.UUID]*elementdomain.Element)}
}
func (r *elementRepo) FindByID(_ context.Context, id uuid.UUID) (*elementdomain.Element, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errors.New("element not found")
	}
	return item, nil
}
func (r *elementRepo) Save(_ context.Context, e *elementdomain.Element) error {
	r.items[e.ID()] = e
	return nil
}

type progressRepo struct {
	item  *progress.Progress
	saved *progress.Progress
}

func newProgressRepo() *progressRepo { return &progressRepo{} }

func (r *progressRepo) FindByEnrollmentID(context.Context, uuid.UUID) (*progress.Progress, error) {
	return r.item, nil
}
func (r *progressRepo) Save(_ context.Context, p *progress.Progress) error {
	r.saved = p
	return nil
}

func mustBlockTitle(t *testing.T, value string) blocktitle.Title {
	t.Helper()
	bt, err := blocktitle.New(value)
	if err != nil {
		t.Fatalf("new block title: %v", err)
	}
	return bt
}
