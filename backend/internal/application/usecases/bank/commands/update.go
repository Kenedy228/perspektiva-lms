package commands

import (
	"context"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
)

type RenameUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewRenameUseCase(r bankports.Repository, audit bankports.AuditRecorder) *RenameUseCase {
	if r == nil {
		panic("bank rename usecase requires repository")
	}
	return &RenameUseCase{r: r, audit: audit}
}

func (uc *RenameUseCase) Execute(ctx context.Context, in RenameInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	b, err := loadBank(ctx, uc.r, in.BankID)
	if err != nil {
		return err
	}

	t, err := buildTitle(in.Title)
	if err != nil {
		return err
	}

	if err := b.Rename(t); err != nil {
		return fmt.Errorf("rename bank: %w", err)
	}

	if err := saveBank(ctx, uc.r, b); err != nil {
		return err
	}

	return recordAudit(ctx, uc.audit, bankports.AuditActionRename, b.ID().String(), in.ActorRole)
}

type AddQuestionsUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewAddQuestionsUseCase(r bankports.Repository, audit bankports.AuditRecorder) *AddQuestionsUseCase {
	if r == nil {
		panic("bank add questions usecase requires repository")
	}
	return &AddQuestionsUseCase{r: r, audit: audit}
}

func (uc *AddQuestionsUseCase) Execute(ctx context.Context, in QuestionIDsInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	b, err := loadBank(ctx, uc.r, in.BankID)
	if err != nil {
		return err
	}

	questionIDs, err := parseRequiredUUIDs(in.QuestionIDs, "question id")
	if err != nil {
		return err
	}

	if err := b.AddQuestions(questionIDs...); err != nil {
		return fmt.Errorf("add bank questions: %w", err)
	}

	if err := saveBank(ctx, uc.r, b); err != nil {
		return err
	}

	return recordAudit(ctx, uc.audit, bankports.AuditActionAddQuestions, b.ID().String(), in.ActorRole)
}

type RemoveQuestionsUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewRemoveQuestionsUseCase(r bankports.Repository, audit bankports.AuditRecorder) *RemoveQuestionsUseCase {
	if r == nil {
		panic("bank remove questions usecase requires repository")
	}
	return &RemoveQuestionsUseCase{r: r, audit: audit}
}

func (uc *RemoveQuestionsUseCase) Execute(ctx context.Context, in QuestionIDsInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	b, err := loadBank(ctx, uc.r, in.BankID)
	if err != nil {
		return err
	}

	questionIDs, err := parseRequiredUUIDs(in.QuestionIDs, "question id")
	if err != nil {
		return err
	}

	if err := b.RemoveQuestions(questionIDs...); err != nil {
		return fmt.Errorf("remove bank questions: %w", err)
	}

	if err := saveBank(ctx, uc.r, b); err != nil {
		return err
	}

	return recordAudit(ctx, uc.audit, bankports.AuditActionRemoveQuestions, b.ID().String(), in.ActorRole)
}

type ClearQuestionsUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewClearQuestionsUseCase(r bankports.Repository, audit bankports.AuditRecorder) *ClearQuestionsUseCase {
	if r == nil {
		panic("bank clear questions usecase requires repository")
	}
	return &ClearQuestionsUseCase{r: r, audit: audit}
}

func (uc *ClearQuestionsUseCase) Execute(ctx context.Context, in BankIDInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	b, err := loadBank(ctx, uc.r, in.BankID)
	if err != nil {
		return err
	}

	b.ClearQuestions()

	if err := saveBank(ctx, uc.r, b); err != nil {
		return err
	}

	return recordAudit(ctx, uc.audit, bankports.AuditActionClearQuestions, b.ID().String(), in.ActorRole)
}

type DeleteUseCase struct {
	r     bankports.Repository
	audit bankports.AuditRecorder
}

func NewDeleteUseCase(r bankports.Repository, audit bankports.AuditRecorder) *DeleteUseCase {
	if r == nil {
		panic("bank delete usecase requires repository")
	}
	return &DeleteUseCase{r: r, audit: audit}
}

func (uc *DeleteUseCase) Execute(ctx context.Context, in BankIDInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	bankID, err := parseRequiredUUID(in.BankID, "bank id")
	if err != nil {
		return err
	}

	if err := uc.r.DeleteByID(ctx, bankID); err != nil {
		return fmt.Errorf("delete bank: %w", err)
	}

	return recordAudit(ctx, uc.audit, bankports.AuditActionDelete, bankID.String(), in.ActorRole)
}
