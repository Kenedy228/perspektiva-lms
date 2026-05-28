package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
	"gitflic.ru/lms/backend/internal/domain/course/block"
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
