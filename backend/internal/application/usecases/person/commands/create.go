package commands

import (
	"context"
	"fmt"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	persondomain "gitflic.ru/lms/backend/internal/domain/person"
	personname "gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type CreateUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewCreateUseCase(r personports.Repository, audit personports.AuditRecorder) *CreateUseCase {
	if r == nil {
		panic("person create usecase requires repository")
	}
	if audit == nil {
		panic("person create usecase requires audit recorder")
	}

	return &CreateUseCase{
		r:     r,
		audit: audit,
	}
}

type CreateInput struct {
	ActorRole  role.Role
	FirstName  string
	LastName   string
	MiddleName string
}

type CreateOutput struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*CreateOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pName, err := personname.New(in.FirstName, in.LastName, in.MiddleName)
	if err != nil {
		return nil, fmt.Errorf("create person name: %w", err)
	}

	p, err := persondomain.New(pName)
	if err != nil {
		return nil, fmt.Errorf("create person aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionCreate, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &CreateOutput{
		ID: p.ID().String(),
	}, nil
}
