package commands

import (
	"context"
	"fmt"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type AssignOrganizationUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewAssignOrganizationUseCase(r personports.Repository, audit personports.AuditRecorder) *AssignOrganizationUseCase {
	if r == nil {
		panic("person assign organization usecase requires repository")
	}
	if audit == nil {
		panic("person assign organization usecase requires audit recorder")
	}

	return &AssignOrganizationUseCase{r: r, audit: audit}
}

type AssignOrganizationInput struct {
	ActorRole      role.Role
	PersonID       string
	OrganizationID string
}

type AssignOrganizationOutput struct {
	ID string
}

func (uc *AssignOrganizationUseCase) Execute(ctx context.Context, in AssignOrganizationInput) (*AssignOrganizationOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	orgID, err := parseRequiredUUID(in.OrganizationID, "organization id")
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.AssignOrganization(orgID); err != nil {
		return nil, fmt.Errorf("assign person organization aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionAssignOrganization, p.ID().String(), in.ActorRole, orgID.String()); err != nil {
		return nil, err
	}

	return &AssignOrganizationOutput{ID: p.ID().String()}, nil
}

type RemoveOrganizationUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewRemoveOrganizationUseCase(r personports.Repository, audit personports.AuditRecorder) *RemoveOrganizationUseCase {
	if r == nil {
		panic("person remove organization usecase requires repository")
	}
	if audit == nil {
		panic("person remove organization usecase requires audit recorder")
	}

	return &RemoveOrganizationUseCase{r: r, audit: audit}
}

type RemoveOrganizationInput struct {
	ActorRole role.Role
	PersonID  string
}

type RemoveOrganizationOutput struct {
	ID string
}

func (uc *RemoveOrganizationUseCase) Execute(ctx context.Context, in RemoveOrganizationInput) (*RemoveOrganizationOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.RemoveOrganization(); err != nil {
		return nil, fmt.Errorf("remove person organization aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionRemoveOrganization, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &RemoveOrganizationOutput{ID: p.ID().String()}, nil
}
