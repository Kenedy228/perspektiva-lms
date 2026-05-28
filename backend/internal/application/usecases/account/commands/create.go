package commands

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type CreateUseCase struct {
	r      accountports.Repository
	hasher accountports.PasswordHasher
	audit  accountports.AuditRecorder
}

func NewCreateUseCase(r accountports.Repository, hasher accountports.PasswordHasher, audit accountports.AuditRecorder) *CreateUseCase {
	if r == nil {
		panic("account create usecase requires repository")
	}
	if hasher == nil {
		panic("account create usecase requires password hasher")
	}
	if audit == nil {
		panic("account create usecase requires audit recorder")
	}

	return &CreateUseCase{r: r, hasher: hasher, audit: audit}
}

type CreateInput struct {
	ActorRole role.Role
	Login     string
	Password  string
	Role      role.Role
	PersonID  string
}

type CreateOutput struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	l, err := login.New(in.Login)
	if err != nil {
		return nil, fmt.Errorf("create account login: %w", err)
	}

	personID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	if err := requireLoginAvailable(ctx, uc.r, l, uuid.Nil); err != nil {
		return nil, err
	}

	if err := requirePersonWithoutAccount(ctx, uc.r, personID); err != nil {
		return nil, err
	}

	hash, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash account password: %w", err)
	}

	acc, err := accountdomain.New(l, hash, in.Role, personID)
	if err != nil {
		return nil, fmt.Errorf("create account aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, acc); err != nil {
		return nil, fmt.Errorf("save account: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, accountports.AuditActionCreate, acc, in.ActorRole); err != nil {
		return nil, err
	}

	return &CreateOutput{ID: acc.ID().String()}, nil
}
