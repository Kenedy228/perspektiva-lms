package commands

import (
	"context"
	"fmt"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	personname "gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type RenameUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewRenameUseCase(r personports.Repository, audit personports.AuditRecorder) *RenameUseCase {
	if r == nil {
		panic("person rename usecase requires repository")
	}
	if audit == nil {
		panic("person rename usecase requires audit recorder")
	}

	return &RenameUseCase{
		r:     r,
		audit: audit,
	}
}

type RenameInput struct {
	ActorRole  role.Role
	PersonID   string
	FirstName  string
	LastName   string
	MiddleName string
}

type RenameOutput struct {
	ID string
}

func (uc *RenameUseCase) Execute(ctx context.Context, in RenameInput) (*RenameOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pName, err := personname.New(in.FirstName, in.LastName, in.MiddleName)
	if err != nil {
		return nil, fmt.Errorf("rename person name: %w", err)
	}

	id, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ChangeName(pName); err != nil {
		return nil, fmt.Errorf("rename person aggregate: %w", err)
	}

	err = uc.r.Save(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionRename, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &RenameOutput{
		ID: p.ID().String(),
	}, nil
}
