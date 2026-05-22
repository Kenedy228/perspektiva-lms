package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

var (
	_ accountports.Repository    = (*AccountRepository)(nil)
	_ accountports.QueryService  = (*AccountRepository)(nil)
	_ accountports.AuditRecorder = (*AccountRepository)(nil)
)

// AccountRepository persists accounts and serves account read models.
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository creates a PostgreSQL account adapter.
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*accountdomain.Account, error) {
	return r.find(ctx, `id = $1`, id)
}

func (r *AccountRepository) FindByLogin(ctx context.Context, l login.Login) (*accountdomain.Account, error) {
	return r.find(ctx, `login = $1`, l.Value())
}

func (r *AccountRepository) LoginExists(ctx context.Context, l login.Login, excludeAccountID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM accounts
			WHERE login = $1 AND ($2::uuid IS NULL OR id <> $2)
		)`, l.Value(), nullUUID(excludeAccountID)).Scan(&exists)
	return exists, err
}

func (r *AccountRepository) PersonHasAccount(ctx context.Context, personID uuid.UUID, excludeAccountID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM accounts
			WHERE person_id = $1 AND ($2::uuid IS NULL OR id <> $2)
		)`, personID, nullUUID(excludeAccountID)).Scan(&exists)
	return exists, err
}

func (r *AccountRepository) Save(ctx context.Context, acc *accountdomain.Account) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO accounts (id, person_id, login, password_hash, role, status, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (id) DO UPDATE SET
			person_id = EXCLUDED.person_id,
			login = EXCLUDED.login,
			password_hash = EXCLUDED.password_hash,
			role = EXCLUDED.role,
			status = EXCLUDED.status,
			updated_at = now()`,
		acc.ID(), acc.PersonID(), acc.Login().Value(), acc.PasswordHash().String(), acc.Role().Kind().String(), acc.Status().String())
	if err != nil {
		return fmt.Errorf("save account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (accountports.AccountView, error) {
	return r.get(ctx, `id = $1`, id)
}

func (r *AccountRepository) GetByPersonID(ctx context.Context, personID uuid.UUID) (accountports.AccountView, error) {
	return r.get(ctx, `person_id = $1`, personID)
}

func (r *AccountRepository) List(ctx context.Context, filter accountports.ListFilter, limit, offset int) ([]accountports.AccountView, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, login, role, person_id::text, status
		FROM accounts
		WHERE ($1 = '' OR role = $1)
			AND ($2 = '' OR status = $2)
			AND ($3 = '' OR lower(login) LIKE lower('%' || $3 || '%'))
		ORDER BY login, id
		LIMIT $4 OFFSET $5`,
		filter.Role.String(), filter.Status.String(), filter.Login, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []accountports.AccountView
	for rows.Next() {
		var view accountports.AccountView
		if err := rows.Scan(&view.ID, &view.Login, &view.Role, &view.PersonID, &view.Status); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}

func (r *AccountRepository) RecordAccountAudit(ctx context.Context, event accountports.AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal account audit event: %w", err)
	}
	var entityID any
	if event.AccountID != "" {
		entityID = event.AccountID
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO audit_events (action, entity_id, actor_role, payload)
		VALUES ($1, $2, $3, $4)`,
		string(event.Action), entityID, event.ActorRole, payload)
	if err != nil {
		return fmt.Errorf("record account audit: %w", err)
	}
	return nil
}

func (r *AccountRepository) find(ctx context.Context, where string, arg any) (*accountdomain.Account, error) {
	var id, personID uuid.UUID
	var loginValue, hashValue, roleValue, statusValue string
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT id, person_id, login, password_hash, role, status
		FROM accounts
		WHERE %s`, where), arg).Scan(&id, &personID, &loginValue, &hashValue, &roleValue, &statusValue)
	if err != nil {
		return nil, err
	}

	l, err := login.New(loginValue)
	if err != nil {
		return nil, fmt.Errorf("restore account login: %w", err)
	}
	h, err := passhash.New(hashValue)
	if err != nil {
		return nil, fmt.Errorf("restore account password hash: %w", err)
	}
	rt, err := role.ParseType(roleValue)
	if err != nil {
		return nil, fmt.Errorf("restore account role: %w", err)
	}
	rv, err := role.New(rt)
	if err != nil {
		return nil, fmt.Errorf("restore account role: %w", err)
	}
	return accountdomain.Restore(id, accountdomain.Params{
		Login:        l,
		PasswordHash: h,
		Role:         rv,
		PersonID:     personID,
		Status:       accountdomain.Status(statusValue),
	})
}

func (r *AccountRepository) get(ctx context.Context, where string, arg any) (accountports.AccountView, error) {
	var view accountports.AccountView
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`
		SELECT id::text, login, role, person_id::text, status
		FROM accounts
		WHERE %s`, where), arg).Scan(&view.ID, &view.Login, &view.Role, &view.PersonID, &view.Status)
	return view, err
}
