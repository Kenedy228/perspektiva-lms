package commands

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ChangeLoginUseCase struct {
	r     accountports.Repository
	audit accountports.AuditRecorder
}

func NewChangeLoginUseCase(r accountports.Repository, audit accountports.AuditRecorder) *ChangeLoginUseCase {
	if r == nil {
		panic("account change login usecase requires repository")
	}
	if audit == nil {
		panic("account change login usecase requires audit recorder")
	}

	return &ChangeLoginUseCase{r: r, audit: audit}
}

type ChangeLoginInput struct {
	ActorRole role.Role
	AccountID string
	Login     string
}

type ChangeLoginOutput struct {
	ID string
}

func (uc *ChangeLoginUseCase) Execute(ctx context.Context, in ChangeLoginInput) (*ChangeLoginOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}

	l, err := login.New(in.Login)
	if err != nil {
		return nil, fmt.Errorf("change account login: %w", err)
	}

	if err := requireLoginAvailable(ctx, uc.r, l, id); err != nil {
		return nil, err
	}

	acc, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find account: %w", err)
	}

	if err := acc.ChangeLogin(l); err != nil {
		return nil, fmt.Errorf("change account login aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionChangeLogin, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &ChangeLoginOutput{ID: acc.ID().String()}, nil
}

type ChangePasswordUseCase struct {
	r      accountports.Repository
	hasher accountports.PasswordHasher
	audit  accountports.AuditRecorder
}

func NewChangePasswordUseCase(r accountports.Repository, hasher accountports.PasswordHasher, audit accountports.AuditRecorder) *ChangePasswordUseCase {
	if r == nil {
		panic("account change password usecase requires repository")
	}
	if hasher == nil {
		panic("account change password usecase requires password hasher")
	}
	if audit == nil {
		panic("account change password usecase requires audit recorder")
	}

	return &ChangePasswordUseCase{r: r, hasher: hasher, audit: audit}
}

type ChangePasswordInput struct {
	ActorRole role.Role
	AccountID string
	Password  string
}

type ChangePasswordOutput struct {
	ID string
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, in ChangePasswordInput) (*ChangePasswordOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}

	acc, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find account: %w", err)
	}

	hash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash account password: %w", err)
	}

	if err := acc.ChangePasswordHash(hash); err != nil {
		return nil, fmt.Errorf("change account password aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionChangePassword, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &ChangePasswordOutput{ID: acc.ID().String()}, nil
}

type ChangeRoleUseCase struct {
	r     accountports.Repository
	audit accountports.AuditRecorder
}

func NewChangeRoleUseCase(r accountports.Repository, audit accountports.AuditRecorder) *ChangeRoleUseCase {
	if r == nil {
		panic("account change role usecase requires repository")
	}
	if audit == nil {
		panic("account change role usecase requires audit recorder")
	}

	return &ChangeRoleUseCase{r: r, audit: audit}
}

type ChangeRoleInput struct {
	ActorRole   role.Role
	AccountID   string
	AccountRole accountdomain.Role
}

type ChangeRoleOutput struct {
	ID string
}

func (uc *ChangeRoleUseCase) Execute(ctx context.Context, in ChangeRoleInput) (*ChangeRoleOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}

	acc, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find account: %w", err)
	}

	if err := acc.ChangeRole(in.AccountRole); err != nil {
		return nil, fmt.Errorf("change account role aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionChangeRole, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &ChangeRoleOutput{ID: acc.ID().String()}, nil
}
