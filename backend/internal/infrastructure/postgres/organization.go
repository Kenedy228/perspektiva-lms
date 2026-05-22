package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	orgdomain "gitflic.ru/lms/backend/internal/domain/organization"
	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	orgname "gitflic.ru/lms/backend/internal/domain/organization/name"
	"github.com/google/uuid"
)

var (
	_ orgports.Repository    = (*OrganizationRepository)(nil)
	_ orgports.QueryService  = (*OrganizationRepository)(nil)
	_ orgports.AuditRecorder = (*OrganizationRepository)(nil)
)

// OrganizationRepository persists organizations and serves organization read models.
type OrganizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository creates a PostgreSQL organization adapter.
func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*orgdomain.Organization, error) {
	var innCode, innType, nameValue string
	err := r.db.QueryRowContext(ctx, `
		SELECT inn, inn_type, name
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&innCode, &innType, &nameValue)
	if err != nil {
		return nil, err
	}

	i, err := inn.New(innCode, inn.Type(innType))
	if err != nil {
		return nil, fmt.Errorf("restore organization inn: %w", err)
	}
	n, err := orgname.New(nameValue)
	if err != nil {
		return nil, fmt.Errorf("restore organization name: %w", err)
	}
	return orgdomain.Restore(id, i, n)
}

func (r *OrganizationRepository) Save(ctx context.Context, o *orgdomain.Organization) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO organizations (id, inn, inn_type, name, updated_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (id) DO UPDATE SET
			inn = EXCLUDED.inn,
			inn_type = EXCLUDED.inn_type,
			name = EXCLUDED.name,
			deleted_at = NULL,
			updated_at = now()`,
		o.ID(), o.INN().Value(), o.INN().Type().String(), o.Name().Value())
	if err != nil {
		return fmt.Errorf("save organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE organizations
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete organization: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) ListByName(ctx context.Context, name string, limit, offset int) ([]orgports.OrganizationShortView, error) {
	return r.list(ctx, "lower(name) LIKE lower('%' || $1 || '%')", name, limit, offset)
}

func (r *OrganizationRepository) ListByINN(ctx context.Context, innValue string, limit, offset int) ([]orgports.OrganizationShortView, error) {
	return r.list(ctx, "inn LIKE $1 || '%'", innValue, limit, offset)
}

func (r *OrganizationRepository) GetDetailsByID(ctx context.Context, id uuid.UUID) (orgports.OrganizationDetailedView, error) {
	var view orgports.OrganizationDetailedView
	var innType string
	err := r.db.QueryRowContext(ctx, `
		SELECT id::text, name, inn, inn_type
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&view.ID, &view.OrganizationName, &view.INN, &innType)
	if err != nil {
		return orgports.OrganizationDetailedView{}, err
	}
	view.INNTitle = inn.Type(innType).Title()
	return view, nil
}

func (r *OrganizationRepository) RecordOrganizationAudit(ctx context.Context, event orgports.AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal organization audit event: %w", err)
	}
	var entityID any
	if event.OrganizationID != "" {
		entityID = event.OrganizationID
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO audit_events (action, entity_id, actor_role, payload)
		VALUES ($1, $2, $3, $4)`,
		string(event.Action), entityID, event.ActorRole, payload)
	if err != nil {
		return fmt.Errorf("record organization audit: %w", err)
	}
	return nil
}

func (r *OrganizationRepository) list(ctx context.Context, where, value string, limit, offset int) ([]orgports.OrganizationShortView, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id::text, name
		FROM organizations
		WHERE deleted_at IS NULL AND %s
		ORDER BY name, id
		LIMIT $2 OFFSET $3`, where), value, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []orgports.OrganizationShortView
	for rows.Next() {
		var view orgports.OrganizationShortView
		if err := rows.Scan(&view.ID, &view.OrganizationName); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}
