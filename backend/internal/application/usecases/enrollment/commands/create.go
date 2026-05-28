package commands

import (
	"context"
	"fmt"
	"time"

	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/application/usecases/enrollment/common"
	enrollmentdomain "gitflic.ru/lms/backend/internal/domain/enrollment"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type CreateUseCase struct {
	repo     enrollmentports.Repository
	progress enrollmentports.ProgressInitializer
}

func NewCreateUseCase(repo enrollmentports.Repository, progress enrollmentports.ProgressInitializer) *CreateUseCase {
	if repo == nil {
		panic("enrollment create usecase requires repository")
	}
	return &CreateUseCase{repo: repo, progress: progress}
}

type CreateInput struct {
	ActorRole     role.Role
	CourseID      string
	AccountID     string
	ActivatedAt   time.Time
	DeactivatedAt time.Time
	Now           time.Time
}

type Output struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}

	courseID, err := parseRequiredUUID(in.CourseID, "course id")
	if err != nil {
		return nil, err
	}
	accountID, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}

	exists, err := uc.repo.ExistsForAccountCourse(ctx, accountID, courseID, uuid.Nil)
	if err != nil {
		return nil, fmt.Errorf("check enrollment uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("%w: аккаунт уже записан на этот курс", common.ErrConflict)
	}

	e, err := enrollmentdomain.NewAt(courseID, accountID, in.ActivatedAt, in.DeactivatedAt, in.Now)
	if err != nil {
		return nil, fmt.Errorf("create enrollment aggregate: %w", err)
	}
	if err := uc.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("save enrollment: %w", err)
	}
	if uc.progress != nil {
		if err := uc.progress.EnsureProgressForEnrollment(ctx, e.ID()); err != nil {
			return nil, fmt.Errorf("ensure enrollment progress: %w", err)
		}
	}

	return &Output{ID: e.ID().String()}, nil
}

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: поле '%s' обязательно", common.ErrInvalidInput, field)
	}
	return id, nil
}
