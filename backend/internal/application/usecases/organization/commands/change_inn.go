package commands

import (
	"context"
	"fmt"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	inn2 "gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type ChangeINNUseCase struct {
	r     orgports.Repository
	audit orgports.AuditRecorder
}

func NewChangeINNUseCase(r orgports.Repository, audit orgports.AuditRecorder) *ChangeINNUseCase {
	if r == nil {
		panic("organization change inn usecase requires repository")
	}
	if audit == nil {
		panic("organization change inn usecase requires audit recorder")
	}

	return &ChangeINNUseCase{
		r:     r,
		audit: audit,
	}
}

type ChangeINNInput struct {
	ActorRole      role.Role
	OrganizationID string
	INN            string
	INNType        string
}

type ChangeINNOutput struct {
	ID string
}

func (uc *ChangeINNUseCase) Execute(ctx context.Context, in ChangeINNInput) (*ChangeINNOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(in.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	i, err := inn2.New(in.INN, inn2.Type(in.INNType))
	if err != nil {
		return nil, fmt.Errorf("create organization inn: %w", err)
	}

	org, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find organization: %w", err)
	}

	if err := org.ChangeINN(i); err != nil {
		return nil, fmt.Errorf("change organization inn: %w", err)
	}

	if err := uc.r.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("save organization: %w", err)
	}

	if err := uc.audit.RecordOrganizationAudit(ctx, orgports.AuditEvent{
		Action:         orgports.AuditActionChangeINN,
		OrganizationID: org.ID().String(),
		ActorRole:      in.ActorRole.Kind().String(),
	}); err != nil {
		return nil, fmt.Errorf("record organization change inn audit: %w", err)
	}

	return &ChangeINNOutput{
		ID: org.ID().String(),
	}, nil
}
