package commands

import (
	"context"
	"fmt"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	orgdomain "gitflic.ru/lms/backend/internal/domain/organization"
	inn2 "gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type CreateUseCase struct {
	r     orgports.Repository
	audit orgports.AuditRecorder
}

func NewCreateUseCase(r orgports.Repository, audit orgports.AuditRecorder) *CreateUseCase {
	if r == nil {
		panic("organization create usecase requires repository")
	}
	if audit == nil {
		panic("organization create usecase requires audit recorder")
	}

	return &CreateUseCase{
		r:     r,
		audit: audit,
	}
}

type CreateInput struct {
	ActorRole role.Role
	INN       string
	INNType   string
	Name      string
}

type CreateOutput struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	i, err := inn2.New(in.INN, inn2.Type(in.INNType))
	if err != nil {
		return nil, fmt.Errorf("create organization inn: %w", err)
	}

	n, err := name.New(in.Name)
	if err != nil {
		return nil, fmt.Errorf("create organization name: %w", err)
	}

	org, err := orgdomain.New(i, n)
	if err != nil {
		return nil, fmt.Errorf("create organization aggregate: %w", err)
	}

	if err = uc.r.Save(ctx, org); err != nil {
		return nil, fmt.Errorf("save organization: %w", err)
	}

	if err := uc.audit.RecordOrganizationAudit(ctx, orgports.AuditEvent{
		Action:         orgports.AuditActionCreate,
		OrganizationID: org.ID().String(),
		ActorRole:      in.ActorRole.Kind().String(),
	}); err != nil {
		return nil, fmt.Errorf("record organization create audit: %w", err)
	}

	return &CreateOutput{
		ID: org.ID().String(),
	}, nil
}
