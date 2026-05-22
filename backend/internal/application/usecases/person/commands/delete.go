package commands

import (
	"context"
	"fmt"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type DeleteByIDUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewDeleteByIDUseCase(r personports.Repository, audit personports.AuditRecorder) *DeleteByIDUseCase {
	if r == nil {
		panic("person delete usecase requires repository")
	}
	if audit == nil {
		panic("person delete usecase requires audit recorder")
	}

	return &DeleteByIDUseCase{
		r:     r,
		audit: audit,
	}
}

type DeleteByIDInput struct {
	ActorRole role.Role
	PersonID  string
}

func (uc *DeleteByIDUseCase) Execute(ctx context.Context, in DeleteByIDInput) error {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return fmt.Errorf("find person: %w", err)
	}

	if err := uc.r.DeleteByID(ctx, pID); err != nil {
		return fmt.Errorf("delete person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionDelete, p.ID().String(), in.ActorRole, ""); err != nil {
		return err
	}

	return nil
}
