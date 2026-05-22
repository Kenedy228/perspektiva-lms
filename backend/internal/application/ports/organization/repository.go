package organization

import (
	"context"

	domainorg "gitflic.ru/lms/backend/internal/domain/organization"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domainorg.Organization, error)
	Save(ctx context.Context, o *domainorg.Organization) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type AuditAction string

const (
	AuditActionCreate    AuditAction = "organization.create"
	AuditActionChangeINN AuditAction = "organization.change_inn"
	AuditActionRename    AuditAction = "organization.rename"
	AuditActionDelete    AuditAction = "organization.delete"
)

type AuditEvent struct {
	Action         AuditAction
	OrganizationID string
	ActorRole      string
}

type AuditRecorder interface {
	RecordOrganizationAudit(ctx context.Context, event AuditEvent) error
}
