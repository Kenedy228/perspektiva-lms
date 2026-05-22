package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/domain/course/title"
	"gitflic.ru/lms/backend/internal/domain/course/version"
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

func TestCreateVersionAndPublishWorkflow(t *testing.T) {
	ctx := context.Background()
	courses := newCourseRepo()
	versions := newVersionRepo()
	blocks := newBlockRepo()

	courseTitle, _ := title.New("Course")
	c, _ := coursedomain.New(courseTitle)
	_ = courses.Save(ctx, c)

	out, err := NewCreateVersionUseCase(courses, versions).Execute(ctx, CreateVersionInput{
		ActorRole: role.NewCreator(),
		CourseID:  c.ID().String(),
		Title:     "v1",
	})
	if err != nil {
		t.Fatalf("create version: %v", err)
	}

	blockOut, err := NewAddBlockUseCase(versions, blocks).Execute(ctx, AddBlockInput{
		ActorRole: role.NewAdmin(),
		VersionID: out.ID,
		Title:     "Block",
	})
	if err != nil {
		t.Fatalf("add block: %v", err)
	}
	if blockOut.ID == "" {
		t.Fatal("expected block id")
	}

	if err := NewPublishVersionUseCase(versions).Execute(ctx, VersionIDInput{
		ActorRole: role.NewAdmin(),
		VersionID: out.ID,
	}); err != nil {
		t.Fatalf("publish version: %v", err)
	}
}

func TestMarkProgressUseCase(t *testing.T) {
	ctx := context.Background()
	repo := newProgressRepo()
	p, _ := progress.New(uuid.New(), uuid.New())
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

type versionRepo struct {
	items map[uuid.UUID]*version.Version
}

func newVersionRepo() *versionRepo { return &versionRepo{items: make(map[uuid.UUID]*version.Version)} }
func (r *versionRepo) FindByID(_ context.Context, id uuid.UUID) (*version.Version, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errors.New("version not found")
	}
	return item, nil
}
func (r *versionRepo) Save(_ context.Context, v *version.Version) error {
	r.items[v.ID()] = v
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
