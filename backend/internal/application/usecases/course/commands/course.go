package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
)

type CreateCourseUseCase struct {
	courses courseports.CourseRepository
}

func NewCreateCourseUseCase(courses courseports.CourseRepository) *CreateCourseUseCase {
	if courses == nil {
		panic("course create usecase requires course repository")
	}
	return &CreateCourseUseCase{courses: courses}
}

func (uc *CreateCourseUseCase) Execute(ctx context.Context, in CreateCourseInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}
	t, err := courseTitle(in.Title)
	if err != nil {
		return nil, err
	}
	c, err := coursedomain.New(t)
	if err != nil {
		return nil, fmt.Errorf("create course aggregate: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("save course: %w", err)
	}
	return &Output{ID: c.ID().String()}, nil
}

type RenameCourseUseCase struct {
	courses courseports.CourseRepository
}

func NewRenameCourseUseCase(courses courseports.CourseRepository) *RenameCourseUseCase {
	if courses == nil {
		panic("course rename usecase requires course repository")
	}
	return &RenameCourseUseCase{courses: courses}
}

func (uc *RenameCourseUseCase) Execute(ctx context.Context, in RenameCourseInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	courseID, err := parseRequiredUUID(in.CourseID, "course id")
	if err != nil {
		return err
	}
	c, err := uc.courses.FindByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("find course: %w", err)
	}
	t, err := courseTitle(in.Title)
	if err != nil {
		return err
	}
	if err := c.ChangeTitle(t); err != nil {
		return fmt.Errorf("rename course: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}
