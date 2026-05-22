package commands

import (
	"context"
	"errors"
	"testing"
	"time"

	enrollmentdomain "gitflic.ru/lms/backend/internal/domain/enrollment"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func TestCreateChecksVersionPolicyUniquenessAndProgress(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)
	repo := &enrollmentRepo{}
	versions := &versionPolicy{allowed: true}
	progress := &progressInitializer{}

	out, err := NewCreateUseCase(repo, versions, progress).Execute(ctx, CreateInput{
		ActorRole:     role.NewAdmin(),
		CourseID:      uuid.New().String(),
		VersionID:     uuid.New().String(),
		AccountID:     uuid.New().String(),
		ActivatedAt:   now,
		DeactivatedAt: now.AddDate(0, 1, 0),
		Now:           now,
	})
	if err != nil {
		t.Fatalf("create enrollment: %v", err)
	}
	if out.ID == "" || repo.saved == nil {
		t.Fatal("expected enrollment to be saved")
	}
	if !progress.called {
		t.Fatal("expected progress initialization")
	}
}

func TestCreateRejectsDeletedVersionAndDuplicateEnrollment(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)

	_, err := NewCreateUseCase(&enrollmentRepo{}, &versionPolicy{allowed: false}, nil).Execute(ctx, CreateInput{
		ActorRole:     role.NewAdmin(),
		CourseID:      uuid.New().String(),
		VersionID:     uuid.New().String(),
		AccountID:     uuid.New().String(),
		ActivatedAt:   now,
		DeactivatedAt: now.AddDate(0, 1, 0),
		Now:           now,
	})
	if err == nil {
		t.Fatal("expected version policy rejection")
	}

	_, err = NewCreateUseCase(&enrollmentRepo{exists: true}, &versionPolicy{allowed: true}, nil).Execute(ctx, CreateInput{
		ActorRole:     role.NewAdmin(),
		CourseID:      uuid.New().String(),
		VersionID:     uuid.New().String(),
		AccountID:     uuid.New().String(),
		ActivatedAt:   now,
		DeactivatedAt: now.AddDate(0, 1, 0),
		Now:           now,
	})
	if err == nil {
		t.Fatal("expected duplicate enrollment rejection")
	}
}

func TestCreateRejectsNonAdmin(t *testing.T) {
	_, err := NewCreateUseCase(&enrollmentRepo{}, &versionPolicy{allowed: true}, nil).Execute(context.Background(), CreateInput{
		ActorRole: role.NewCreator(),
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

type enrollmentRepo struct {
	saved  *enrollmentdomain.Enrollment
	exists bool
}

func (r *enrollmentRepo) FindByID(context.Context, uuid.UUID) (*enrollmentdomain.Enrollment, error) {
	return nil, errors.New("not found")
}

func (r *enrollmentRepo) Save(_ context.Context, e *enrollmentdomain.Enrollment) error {
	r.saved = e
	return nil
}

func (r *enrollmentRepo) ExistsForAccountCourseVersion(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID) (bool, error) {
	return r.exists, nil
}

type versionPolicy struct {
	allowed bool
}

func (p *versionPolicy) CanEnrollVersion(context.Context, uuid.UUID) (bool, error) {
	return p.allowed, nil
}

type progressInitializer struct {
	called bool
}

func (p *progressInitializer) EnsureProgressForEnrollment(context.Context, uuid.UUID, uuid.UUID) error {
	p.called = true
	return nil
}
