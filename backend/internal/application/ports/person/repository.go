package person

import (
	"context"

	"gitflic.ru/lms/backend/internal/domain/person"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*person.Person, error)
	SNILSExists(ctx context.Context, value snils.SNILS, excludePersonID uuid.UUID) (bool, error)
	Save(ctx context.Context, p *person.Person) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type AuditAction string

const (
	AuditActionCreate             AuditAction = "person.create"
	AuditActionCreateWithProfile  AuditAction = "person.create_with_profile"
	AuditActionRename             AuditAction = "person.rename"
	AuditActionAttachProfile      AuditAction = "person.attach_profile"
	AuditActionReplaceProfile     AuditAction = "person.replace_profile"
	AuditActionDetachProfile      AuditAction = "person.detach_profile"
	AuditActionChangeSNILS        AuditAction = "person.change_snils"
	AuditActionChangeDateOfBirth  AuditAction = "person.change_date_of_birth"
	AuditActionChangeJobTitle     AuditAction = "person.change_job_title"
	AuditActionChangeEducation    AuditAction = "person.change_education"
	AuditActionAssignOrganization AuditAction = "person.assign_organization"
	AuditActionRemoveOrganization AuditAction = "person.remove_organization"
	AuditActionDelete             AuditAction = "person.delete"
)

type AuditEvent struct {
	Action         AuditAction
	PersonID       string
	ActorRole      string
	OrganizationID string
}

type AuditRecorder interface {
	RecordPersonAudit(ctx context.Context, event AuditEvent) error
}
