package commands

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
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
