package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	"gitflic.ru/lms/backend/internal/domain/course/version"
)

type CreateVersionUseCase struct {
	courses  courseports.CourseRepository
	versions courseports.VersionRepository
}

func NewCreateVersionUseCase(courses courseports.CourseRepository, versions courseports.VersionRepository) *CreateVersionUseCase {
	if courses == nil || versions == nil {
		panic("course create version usecase requires repositories")
	}
	return &CreateVersionUseCase{courses: courses, versions: versions}
}

func (uc *CreateVersionUseCase) Execute(ctx context.Context, in CreateVersionInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}
	courseID, err := parseRequiredUUID(in.CourseID, "course id")
	if err != nil {
		return nil, err
	}
	c, err := uc.courses.FindByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("find course: %w", err)
	}
	t, err := versionTitle(in.Title)
	if err != nil {
		return nil, err
	}
	v, err := version.New(t)
	if err != nil {
		return nil, fmt.Errorf("create version aggregate: %w", err)
	}
	if err := c.AddVersionID(v.ID()); err != nil {
		return nil, fmt.Errorf("attach version to course: %w", err)
	}
	if err := uc.versions.Save(ctx, v); err != nil {
		return nil, fmt.Errorf("save version: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("save course: %w", err)
	}
	return &Output{ID: v.ID().String()}, nil
}

type AddBlockUseCase struct {
	versions courseports.VersionRepository
	blocks   courseports.BlockRepository
}

func NewAddBlockUseCase(versions courseports.VersionRepository, blocks courseports.BlockRepository) *AddBlockUseCase {
	if versions == nil || blocks == nil {
		panic("course add block usecase requires repositories")
	}
	return &AddBlockUseCase{versions: versions, blocks: blocks}
}

func (uc *AddBlockUseCase) Execute(ctx context.Context, in AddBlockInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}
	versionID, err := parseRequiredUUID(in.VersionID, "version id")
	if err != nil {
		return nil, err
	}
	v, err := uc.versions.FindByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("find version: %w", err)
	}
	t, err := blockTitle(in.Title)
	if err != nil {
		return nil, err
	}
	b, err := block.New(t)
	if err != nil {
		return nil, fmt.Errorf("create block aggregate: %w", err)
	}
	if err := v.AddBlockID(b.ID()); err != nil {
		return nil, fmt.Errorf("attach block to version: %w", err)
	}
	if err := uc.blocks.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save block: %w", err)
	}
	if err := uc.versions.Save(ctx, v); err != nil {
		return nil, fmt.Errorf("save version: %w", err)
	}
	return &Output{ID: b.ID().String()}, nil
}

type PublishVersionUseCase struct {
	versions courseports.VersionRepository
}

func NewPublishVersionUseCase(versions courseports.VersionRepository) *PublishVersionUseCase {
	if versions == nil {
		panic("course publish version usecase requires version repository")
	}
	return &PublishVersionUseCase{versions: versions}
}

func (uc *PublishVersionUseCase) Execute(ctx context.Context, in VersionIDInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	versionID, err := parseRequiredUUID(in.VersionID, "version id")
	if err != nil {
		return err
	}
	v, err := uc.versions.FindByID(ctx, versionID)
	if err != nil {
		return fmt.Errorf("find version: %w", err)
	}
	if err := v.Publish(); err != nil {
		return fmt.Errorf("publish version: %w", err)
	}
	if err := uc.versions.Save(ctx, v); err != nil {
		return fmt.Errorf("save version: %w", err)
	}
	return nil
}
