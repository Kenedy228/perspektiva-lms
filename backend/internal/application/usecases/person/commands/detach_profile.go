package commands

import (
	"context"
	"fmt"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type DetachProfileUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewDetachProfileUseCase(r personports.Repository, audit personports.AuditRecorder) *DetachProfileUseCase {
	if r == nil {
		panic("person detach profile usecase requires repository")
	}
	if audit == nil {
		panic("person detach profile usecase requires audit recorder")
	}

	return &DetachProfileUseCase{
		r:     r,
		audit: audit,
	}
}

type DetachProfileInput struct {
	ActorRole role.Role
	PersonID  string
}

type DetachProfileOutput struct {
	PersonID string
}

func (uc *DetachProfileUseCase) Execute(ctx context.Context, in DetachProfileInput) (*DetachProfileOutput, error) {
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

	if p.HasProfile() {
		p.DetachProfile()
		if err := uc.r.Save(ctx, p); err != nil {
			return nil, fmt.Errorf("save person: %w", err)
		}

		if err := recordAudit(ctx, uc.audit, personports.AuditActionDetachProfile, p.ID().String(), in.ActorRole, ""); err != nil {
			return nil, err
		}
	}

	return &DetachProfileOutput{
		PersonID: p.ID().String(),
	}, nil
}
