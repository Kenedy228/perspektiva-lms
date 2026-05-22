package commands

import (
	"context"
	"fmt"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type DeleteByIDUseCase struct {
	r     orgports.Repository
	audit orgports.AuditRecorder
}

func NewDeleteByIDUseCase(r orgports.Repository, audit orgports.AuditRecorder) *DeleteByIDUseCase {
	if r == nil {
		panic("organization delete usecase requires repository")
	}
	if audit == nil {
		panic("organization delete usecase requires audit recorder")
	}

	return &DeleteByIDUseCase{
		r:     r,
		audit: audit,
	}
}

type DeleteByIDInput struct {
	ActorRole      role.Role
	OrganizationID string
}

func (uc *DeleteByIDUseCase) Execute(ctx context.Context, in DeleteByIDInput) error {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return err
	}

	id, err := uuid.Parse(in.OrganizationID)
	if err != nil {
		return fmt.Errorf("parse organization id: %w", err)
	}

	org, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find organization: %w", err)
	}

	if err := uc.r.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("delete organization: %w", err)
	}

	if err := uc.audit.RecordOrganizationAudit(ctx, orgports.AuditEvent{
		Action:         orgports.AuditActionDelete,
		OrganizationID: org.ID().String(),
		ActorRole:      in.ActorRole.Kind().String(),
	}); err != nil {
		return fmt.Errorf("record organization delete audit: %w", err)
	}

	return nil
}
