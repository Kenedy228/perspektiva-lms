package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"github.com/google/uuid"
)

type MarkProgressUseCase struct {
	progress courseports.ProgressRepository
	orgScope enrollmentports.OrganizationScope
}

func NewMarkProgressUseCase(progress courseports.ProgressRepository, orgScope enrollmentports.OrganizationScope) *MarkProgressUseCase {
	if progress == nil {
		panic("course mark progress usecase requires progress repository")
	}
	return &MarkProgressUseCase{progress: progress, orgScope: orgScope}
}

func (uc *MarkProgressUseCase) Execute(_ context.Context, _ MarkProgressInput) error {
	return fmt.Errorf("%w: изменение прогресса не разрешено", common.ErrForbidden)
}

type UnmarkElementCompletedUseCase struct {
	progress courseports.ProgressRepository
	orgScope enrollmentports.OrganizationScope
}

func NewUnmarkElementCompletedUseCase(progress courseports.ProgressRepository, orgScope enrollmentports.OrganizationScope) *UnmarkElementCompletedUseCase {
	if progress == nil {
		panic("course unmark element usecase requires progress repository")
	}
	return &UnmarkElementCompletedUseCase{progress: progress, orgScope: orgScope}
}

func (uc *UnmarkElementCompletedUseCase) Execute(_ context.Context, _ UnmarkElementCompletedInput) error {
	return fmt.Errorf("%w: изменение прогресса не разрешено", common.ErrForbidden)
}

type GetProgressUseCase struct {
	progress courseports.ProgressRepository
	orgScope enrollmentports.OrganizationScope
}

func NewGetProgressUseCase(progress courseports.ProgressRepository, orgScope enrollmentports.OrganizationScope) *GetProgressUseCase {
	if progress == nil {
		panic("course get progress usecase requires progress repository")
	}
	return &GetProgressUseCase{progress: progress, orgScope: orgScope}
}

func (uc *GetProgressUseCase) Execute(ctx context.Context, in GetProgressInput) (*ProgressOutput, error) {
	if err := common.RequireProgressAccess(in.ActorRole); err != nil {
		return nil, err
	}
	if err := common.RequireOrganizationScope(ctx, uc.orgScope, in.ActorRole, in.ActorPersonID, in.EnrollmentID); err != nil {
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
		Percent:             p.CompletionPercent(in.TotalTrackedItems),
		TotalTrackedItems:   in.TotalTrackedItems,
		CompletedElementIDs: elementIDs,
	}, nil
}
