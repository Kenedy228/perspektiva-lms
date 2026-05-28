package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	elementdomain "gitflic.ru/lms/backend/internal/domain/course/element"
)

type AddBlockToCourseUseCase struct {
	courses courseports.CourseRepository
	blocks  courseports.BlockRepository
}

func NewAddBlockToCourseUseCase(courses courseports.CourseRepository, blocks courseports.BlockRepository) *AddBlockToCourseUseCase {
	if courses == nil || blocks == nil {
		panic("course add block usecase requires repositories")
	}
	return &AddBlockToCourseUseCase{courses: courses, blocks: blocks}
}

func (uc *AddBlockToCourseUseCase) Execute(ctx context.Context, in AddBlockToCourseInput) (*Output, error) {
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
	t, err := blockTitle(in.Title)
	if err != nil {
		return nil, err
	}
	b, err := block.New(t)
	if err != nil {
		return nil, fmt.Errorf("create block aggregate: %w", err)
	}
	if err := c.AddBlockID(b.ID()); err != nil {
		return nil, fmt.Errorf("attach block to course: %w", err)
	}
	if err := uc.blocks.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save block: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return nil, fmt.Errorf("save course: %w", err)
	}
	return &Output{ID: b.ID().String()}, nil
}

type RemoveBlockFromCourseUseCase struct {
	courses courseports.CourseRepository
}

func NewRemoveBlockFromCourseUseCase(courses courseports.CourseRepository) *RemoveBlockFromCourseUseCase {
	if courses == nil {
		panic("course remove block usecase requires course repository")
	}
	return &RemoveBlockFromCourseUseCase{courses: courses}
}

func (uc *RemoveBlockFromCourseUseCase) Execute(ctx context.Context, in RemoveBlockFromCourseInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	courseID, err := parseRequiredUUID(in.CourseID, "course id")
	if err != nil {
		return err
	}
	blockID, err := parseRequiredUUID(in.BlockID, "block id")
	if err != nil {
		return err
	}
	c, err := uc.courses.FindByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("find course: %w", err)
	}
	if err := c.RemoveBlockID(blockID); err != nil {
		return fmt.Errorf("remove course block: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

type MoveCourseBlockUseCase struct {
	courses courseports.CourseRepository
}

func NewMoveCourseBlockUseCase(courses courseports.CourseRepository) *MoveCourseBlockUseCase {
	if courses == nil {
		panic("course move block usecase requires course repository")
	}
	return &MoveCourseBlockUseCase{courses: courses}
}

func (uc *MoveCourseBlockUseCase) Execute(ctx context.Context, in MoveCourseBlockInput) error {
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
	if err := c.MoveBlock(in.From, in.To); err != nil {
		return fmt.Errorf("move course block: %w", err)
	}
	if err := uc.courses.Save(ctx, c); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	return nil
}

type AddElementToBlockUseCase struct {
	blocks   courseports.BlockRepository
	elements courseports.ElementRepository
}

func NewAddElementToBlockUseCase(blocks courseports.BlockRepository, elements courseports.ElementRepository) *AddElementToBlockUseCase {
	if blocks == nil || elements == nil {
		panic("course add element usecase requires repositories")
	}
	return &AddElementToBlockUseCase{blocks: blocks, elements: elements}
}

func (uc *AddElementToBlockUseCase) Execute(ctx context.Context, in AddElementToBlockInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}
	blockID, err := parseRequiredUUID(in.BlockID, "block id")
	if err != nil {
		return nil, err
	}
	b, err := uc.blocks.FindByID(ctx, blockID)
	if err != nil {
		return nil, fmt.Errorf("find block: %w", err)
	}
	t, err := elementTitle(in.Title)
	if err != nil {
		return nil, err
	}
	content, err := buildElementContent(in.Content)
	if err != nil {
		return nil, err
	}
	e, err := elementdomain.New(t, content)
	if err != nil {
		return nil, fmt.Errorf("create element aggregate: %w", err)
	}
	mode, err := completionMode(in.Content.CompletionMode)
	if err != nil {
		return nil, err
	}
	if err := e.ChangeCompletionMode(mode); err != nil {
		return nil, fmt.Errorf("change completion mode: %w", err)
	}
	if err := b.AddElementID(e.ID()); err != nil {
		return nil, fmt.Errorf("attach element to block: %w", err)
	}
	if err := uc.elements.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("save element: %w", err)
	}
	if err := uc.blocks.Save(ctx, b); err != nil {
		return nil, fmt.Errorf("save block: %w", err)
	}
	return &Output{ID: e.ID().String()}, nil
}

type RemoveElementFromBlockUseCase struct {
	blocks courseports.BlockRepository
}

func NewRemoveElementFromBlockUseCase(blocks courseports.BlockRepository) *RemoveElementFromBlockUseCase {
	if blocks == nil {
		panic("course remove element usecase requires block repository")
	}
	return &RemoveElementFromBlockUseCase{blocks: blocks}
}

func (uc *RemoveElementFromBlockUseCase) Execute(ctx context.Context, in RemoveElementFromBlockInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	blockID, err := parseRequiredUUID(in.BlockID, "block id")
	if err != nil {
		return err
	}
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return err
	}
	b, err := uc.blocks.FindByID(ctx, blockID)
	if err != nil {
		return fmt.Errorf("find block: %w", err)
	}
	if err := b.RemoveElementID(elementID); err != nil {
		return fmt.Errorf("remove block element: %w", err)
	}
	if err := uc.blocks.Save(ctx, b); err != nil {
		return fmt.Errorf("save block: %w", err)
	}
	return nil
}

type MoveBlockElementUseCase struct {
	blocks courseports.BlockRepository
}

func NewMoveBlockElementUseCase(blocks courseports.BlockRepository) *MoveBlockElementUseCase {
	if blocks == nil {
		panic("course move block element usecase requires block repository")
	}
	return &MoveBlockElementUseCase{blocks: blocks}
}

func (uc *MoveBlockElementUseCase) Execute(ctx context.Context, in MoveBlockElementInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	blockID, err := parseRequiredUUID(in.BlockID, "block id")
	if err != nil {
		return err
	}
	b, err := uc.blocks.FindByID(ctx, blockID)
	if err != nil {
		return fmt.Errorf("find block: %w", err)
	}
	if err := b.MoveElement(in.From, in.To); err != nil {
		return fmt.Errorf("move block element: %w", err)
	}
	if err := uc.blocks.Save(ctx, b); err != nil {
		return fmt.Errorf("save block: %w", err)
	}
	return nil
}

type ChangeElementCompletionModeUseCase struct {
	elements courseports.ElementRepository
}

func NewChangeElementCompletionModeUseCase(elements courseports.ElementRepository) *ChangeElementCompletionModeUseCase {
	if elements == nil {
		panic("course change element completion mode usecase requires element repository")
	}
	return &ChangeElementCompletionModeUseCase{elements: elements}
}

func (uc *ChangeElementCompletionModeUseCase) Execute(ctx context.Context, in ChangeElementCompletionModeInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return err
	}
	mode := elementdomain.CompletionMode(in.CompletionMode)
	switch mode {
	case elementdomain.CompletionModeNone, elementdomain.CompletionModeManual:
	default:
		return fmt.Errorf("%w: неизвестный режим отслеживания %q", common.ErrInvalidInput, in.CompletionMode)
	}
	e, err := uc.elements.FindByID(ctx, elementID)
	if err != nil {
		return fmt.Errorf("find element: %w", err)
	}
	if err := e.ChangeCompletionMode(mode); err != nil {
		return fmt.Errorf("change element completion mode: %w", err)
	}
	if err := uc.elements.Save(ctx, e); err != nil {
		return fmt.Errorf("save element: %w", err)
	}
	return nil
}
