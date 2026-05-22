package commands

import (
	"context"
	"fmt"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type RenameUseCase struct {
	r     orgports.Repository
	audit orgports.AuditRecorder
}

func NewRenameUseCase(r orgports.Repository, audit orgports.AuditRecorder) *RenameUseCase {
	if r == nil {
		panic("organization rename usecase requires repository")
	}
	if audit == nil {
		panic("organization rename usecase requires audit recorder")
	}

	return &RenameUseCase{
		r:     r,
		audit: audit,
	}
}

type RenameInput struct {
	ActorRole      role.Role
	OrganizationID string
	Name           string
}

type RenameOutput struct {
	ID string
}

func (uc *RenameUseCase) Execute(ctx context.Context, in RenameInput) (*RenameOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(in.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	n, err := name.New(in.Name)
	if err != nil {
		return nil, fmt.Errorf("create organization name: %w", err)
	}

	org, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find organization: %w", err)
	}

	if err := org.ChangeName(n); err != nil {
		return nil, fmt.Errorf("rename organization: %w", err)
	}

	if err := uc.r.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("save organization: %w", err)
	}

	if err := uc.audit.RecordOrganizationAudit(ctx, orgports.AuditEvent{
		Action:         orgports.AuditActionRename,
		OrganizationID: org.ID().String(),
		ActorRole:      in.ActorRole.Kind().String(),
	}); err != nil {
		return nil, fmt.Errorf("record organization rename audit: %w", err)
	}

	return &RenameOutput{
		ID: org.ID().String(),
	}, nil
}
