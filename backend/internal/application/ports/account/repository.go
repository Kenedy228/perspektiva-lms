package account

import (
	"context"

	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error)
	FindByLogin(ctx context.Context, login login.Login) (*accountdomain.Account, error)
	LoginExists(ctx context.Context, login login.Login, excludeAccountID uuid.UUID) (bool, error)
	PersonHasAccount(ctx context.Context, personID uuid.UUID, excludeAccountID uuid.UUID) (bool, error)
	Save(ctx context.Context, acc *accountdomain.Account) error
}

type PasswordHasher interface {
	Hash(plain string) (passhash.Hash, error)
}

type AuditAction string

const (
	AuditActionCreate         AuditAction = "account.create"
	AuditActionChangeLogin    AuditAction = "account.change_login"
	AuditActionChangePassword AuditAction = "account.change_password"
	AuditActionChangeRole     AuditAction = "account.change_role"
	AuditActionBlock          AuditAction = "account.block"
	AuditActionActivate       AuditAction = "account.activate"
	AuditActionDelete         AuditAction = "account.delete"
)

type AuditEvent struct {
	Action    AuditAction
	AccountID string
	PersonID  string
	ActorRole string
}

type AuditRecorder interface {
	RecordAccountAudit(ctx context.Context, event AuditEvent) error
}
