package commands

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type BlockUseCase struct {
	r     accountports.Repository
	audit accountports.AuditRecorder
}

func NewBlockUseCase(r accountports.Repository, audit accountports.AuditRecorder) *BlockUseCase {
	if r == nil {
		panic("account block usecase requires repository")
	}
	if audit == nil {
		panic("account block usecase requires audit recorder")
	}

	return &BlockUseCase{r: r, audit: audit}
}

type BlockInput struct {
	ActorRole role.Role
	AccountID string
}

type BlockOutput struct {
	ID string
}

func (uc *BlockUseCase) Execute(ctx context.Context, in BlockInput) (*BlockOutput, error) {
	acc, err := loadAdminAccount(ctx, uc.r, in.ActorRole, in.AccountID)
	if err != nil {
		return nil, err
	}

	if err := acc.Block(); err != nil {
		return nil, fmt.Errorf("block account aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionBlock, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &BlockOutput{ID: acc.ID().String()}, nil
}

type ActivateUseCase struct {
	r     accountports.Repository
	audit accountports.AuditRecorder
}

func NewActivateUseCase(r accountports.Repository, audit accountports.AuditRecorder) *ActivateUseCase {
	if r == nil {
		panic("account activate usecase requires repository")
	}
	if audit == nil {
		panic("account activate usecase requires audit recorder")
	}

	return &ActivateUseCase{r: r, audit: audit}
}

type ActivateInput struct {
	ActorRole role.Role
	AccountID string
}

type ActivateOutput struct {
	ID string
}

func (uc *ActivateUseCase) Execute(ctx context.Context, in ActivateInput) (*ActivateOutput, error) {
	acc, err := loadAdminAccount(ctx, uc.r, in.ActorRole, in.AccountID)
	if err != nil {
		return nil, err
	}

	if err := acc.Activate(); err != nil {
		return nil, fmt.Errorf("activate account aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionActivate, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &ActivateOutput{ID: acc.ID().String()}, nil
}

type DeleteUseCase struct {
	r     accountports.Repository
	audit accountports.AuditRecorder
}

func NewDeleteUseCase(r accountports.Repository, audit accountports.AuditRecorder) *DeleteUseCase {
	if r == nil {
		panic("account delete usecase requires repository")
	}
	if audit == nil {
		panic("account delete usecase requires audit recorder")
	}

	return &DeleteUseCase{r: r, audit: audit}
}

type DeleteInput struct {
	ActorRole role.Role
	AccountID string
}

func (uc *DeleteUseCase) Execute(ctx context.Context, in DeleteInput) error {
	acc, err := loadAdminAccount(ctx, uc.r, in.ActorRole, in.AccountID)
	if err != nil {
		return err
	}

	if err := acc.Delete(); err != nil {
		return fmt.Errorf("delete account aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionDelete, acc, in.ActorRole); err != nil {
		return err
	}

	return nil
}

func loadAdminAccount(ctx context.Context, r accountports.Repository, actor role.Role, accountID string) (*accountdomain.Account, error) {
	if err := common.RequireAdmin(actor); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(accountID, "account id")
	if err != nil {
		return nil, err
	}

	acc, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find account: %w", err)
	}

	return acc, nil
}
