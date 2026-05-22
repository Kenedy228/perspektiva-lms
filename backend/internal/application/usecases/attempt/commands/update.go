package commands

import (
	"context"
	"fmt"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	"gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
)

type AddAnswerUseCase struct {
	r attemptports.Repository
}

func NewAddAnswerUseCase(r attemptports.Repository) *AddAnswerUseCase {
	if r == nil {
		panic("attempt add answer usecase requires repository")
	}
	return &AddAnswerUseCase{r: r}
}

func (uc *AddAnswerUseCase) Execute(ctx context.Context, in AddAnswerInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}

	a, err := loadAttempt(ctx, uc.r, in.AttemptID)
	if err != nil {
		return err
	}

	questionID, err := parseRequiredUUID(in.QuestionID, "question id")
	if err != nil {
		return err
	}

	if err := a.AddAnswer(questionID, in.Answer, in.AnsweredAt); err != nil {
		return fmt.Errorf("add attempt answer: %w", err)
	}

	return saveAttempt(ctx, uc.r, a)
}

type FinishUseCase struct {
	r attemptports.Repository
}

func NewFinishUseCase(r attemptports.Repository) *FinishUseCase {
	if r == nil {
		panic("attempt finish usecase requires repository")
	}
	return &FinishUseCase{r: r}
}

func (uc *FinishUseCase) Execute(ctx context.Context, in FinishInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}

	a, err := loadAttempt(ctx, uc.r, in.AttemptID)
	if err != nil {
		return err
	}

	if err := a.Finish(in.FinishedAt); err != nil {
		return fmt.Errorf("finish attempt: %w", err)
	}

	return saveAttempt(ctx, uc.r, a)
}

type ExpireUseCase struct {
	r attemptports.Repository
}

func NewExpireUseCase(r attemptports.Repository) *ExpireUseCase {
	if r == nil {
		panic("attempt expire usecase requires repository")
	}
	return &ExpireUseCase{r: r}
}

func (uc *ExpireUseCase) Execute(ctx context.Context, in ExpireInput) error {
	a, err := loadAttempt(ctx, uc.r, in.AttemptID)
	if err != nil {
		return err
	}

	if err := a.SetExpired(in.ExpiredAt); err != nil {
		return fmt.Errorf("expire attempt: %w", err)
	}

	return saveAttempt(ctx, uc.r, a)
}

type CancelUseCase struct {
	r attemptports.Repository
}

func NewCancelUseCase(r attemptports.Repository) *CancelUseCase {
	if r == nil {
		panic("attempt cancel usecase requires repository")
	}
	return &CancelUseCase{r: r}
}

func (uc *CancelUseCase) Execute(ctx context.Context, in CancelInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}

	a, err := loadAttempt(ctx, uc.r, in.AttemptID)
	if err != nil {
		return err
	}

	if err := a.Cancel(in.CancelledAt); err != nil {
		return fmt.Errorf("cancel attempt: %w", err)
	}

	return saveAttempt(ctx, uc.r, a)
}

type InterruptUseCase struct {
	r attemptports.Repository
}

func NewInterruptUseCase(r attemptports.Repository) *InterruptUseCase {
	if r == nil {
		panic("attempt interrupt usecase requires repository")
	}
	return &InterruptUseCase{r: r}
}

func (uc *InterruptUseCase) Execute(ctx context.Context, in InterruptInput) error {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return err
	}

	a, err := loadAttempt(ctx, uc.r, in.AttemptID)
	if err != nil {
		return err
	}

	if err := a.Interrupt(in.InterruptedAt); err != nil {
		return fmt.Errorf("interrupt attempt: %w", err)
	}

	return saveAttempt(ctx, uc.r, a)
}
