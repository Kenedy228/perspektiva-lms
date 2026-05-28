package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"github.com/google/uuid"
)

type MarkProgressUseCase struct {
	progress courseports.ProgressRepository
}

func NewMarkProgressUseCase(progress courseports.ProgressRepository) *MarkProgressUseCase {
	if progress == nil {
		panic("course mark progress usecase requires progress repository")
	}
	return &MarkProgressUseCase{progress: progress}
}

func (uc *MarkProgressUseCase) Execute(ctx context.Context, in MarkProgressInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}
	enrollmentID, err := parseRequiredUUID(in.EnrollmentID, "enrollment id")
	if err != nil {
		return err
	}
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return err
	}
	p, err := uc.progress.FindByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return fmt.Errorf("find progress: %w", err)
	}
	if err := p.MarkElement(elementID, in.MarkerType, in.At); err != nil {
		return fmt.Errorf("mark course progress: %w", err)
	}
	if err := uc.progress.Save(ctx, p); err != nil {
		return fmt.Errorf("save progress: %w", err)
	}
	return nil
}

type UnmarkElementCompletedUseCase struct {
	progress courseports.ProgressRepository
}

func NewUnmarkElementCompletedUseCase(progress courseports.ProgressRepository) *UnmarkElementCompletedUseCase {
	if progress == nil {
		panic("course unmark element usecase requires progress repository")
	}
	return &UnmarkElementCompletedUseCase{progress: progress}
}

func (uc *UnmarkElementCompletedUseCase) Execute(ctx context.Context, in UnmarkElementCompletedInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}
	enrollmentID, err := parseRequiredUUID(in.EnrollmentID, "enrollment id")
	if err != nil {
		return err
	}
	elementID, err := parseRequiredUUID(in.ElementID, "element id")
	if err != nil {
		return err
	}
	p, err := uc.progress.FindByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return fmt.Errorf("find progress: %w", err)
	}
	if err := p.UnmarkCompleted(elementID); err != nil {
		return fmt.Errorf("unmark element: %w", err)
	}
	if err := uc.progress.Save(ctx, p); err != nil {
		return fmt.Errorf("save progress: %w", err)
	}
	return nil
}

type GetProgressUseCase struct {
	progress courseports.ProgressRepository
}

func NewGetProgressUseCase(progress courseports.ProgressRepository) *GetProgressUseCase {
	if progress == nil {
		panic("course get progress usecase requires progress repository")
	}
	return &GetProgressUseCase{progress: progress}
}

func (uc *GetProgressUseCase) Execute(ctx context.Context, in GetProgressInput) (*ProgressOutput, error) {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return nil, err
	}
	enrollmentID, err := parseRequiredUUID(in.EnrollmentID, "enrollment id")
	if err != nil {
		return nil, err
	}
	p, err := uc.progress.FindByEnrollmentID(ctx, enrollmentID)
	if err != nil {
		return nil, fmt.Errorf("find progress: %w", err)
	}
	markers := p.Markers()
	elementIDs := make([]uuid.UUID, 0, len(markers))
	for elementID := range markers {
		elementIDs = append(elementIDs, elementID)
	}
	return &ProgressOutput{
		CompletedCount:      p.CompletedCount(),
		Percent:             p.CompletionPercent(len(markers)),
		CompletedElementIDs: elementIDs,
	}, nil
}
