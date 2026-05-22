package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	persondomain "gitflic.ru/lms/backend/internal/domain/person"
	personname "gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

var (
	_ personports.Repository    = (*PersonRepository)(nil)
	_ personports.QueryService  = (*PersonRepository)(nil)
	_ personports.AuditRecorder = (*PersonRepository)(nil)
)

// PersonRepository persists persons and serves person read models.
type PersonRepository struct {
	db *sql.DB
}

// NewPersonRepository creates a PostgreSQL person adapter.
func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) FindByID(ctx context.Context, id uuid.UUID) (*persondomain.Person, error) {
	var firstName, lastName, middleName, jobTitleValue, educationValue string
	var snilsValue sql.NullString
	var dobValue sql.NullTime
	var organizationID uuid.NullUUID
	err := r.db.QueryRowContext(ctx, `
		SELECT first_name, last_name, middle_name, snils, date_of_birth, job_title, education, organization_id
		FROM persons
		WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&firstName, &lastName, &middleName, &snilsValue, &dobValue, &jobTitleValue, &educationValue, &organizationID)
	if err != nil {
		return nil, err
	}

	n, err := personname.New(firstName, lastName, middleName)
	if err != nil {
		return nil, fmt.Errorf("restore person name: %w", err)
	}

	var prof *profile.Profile
	if snilsValue.Valid && dobValue.Valid {
		s, err := snils.New(snilsValue.String)
		if err != nil {
			return nil, fmt.Errorf("restore person snils: %w", err)
		}
		db, err := dob.New(dobValue.Time, time.Now())
		if err != nil {
			return nil, fmt.Errorf("restore person date of birth: %w", err)
		}
		jt, err := jobtitle.New(jobTitleValue)
		if err != nil {
			return nil, fmt.Errorf("restore person job title: %w", err)
		}
		ed, err := education.New(educationValue)
		if err != nil {
			return nil, fmt.Errorf("restore person education: %w", err)
		}
		params := profile.Params{Snils: s, DateOfBirth: db, JobTitle: jt, Education: ed}
		if organizationID.Valid {
			params.OrganizationID = organizationID.UUID
		}
		restored, err := profile.New(params)
		if err != nil {
			return nil, fmt.Errorf("restore person profile: %w", err)
		}
		prof = &restored
	}

	return persondomain.Restore(id, n, prof)
}

func (r *PersonRepository) SNILSExists(ctx context.Context, value snils.SNILS, excludePersonID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM persons
			WHERE snils = $1 AND deleted_at IS NULL AND ($2::uuid IS NULL OR id <> $2)
		)`, value.Value(), nullUUID(excludePersonID)).Scan(&exists)
	return exists, err
}

func (r *PersonRepository) Save(ctx context.Context, p *persondomain.Person) error {
	prof, hasProfile := p.Profile()
	var snilsValue, dobValue, organizationID, jobTitleValue, educationValue any
	if hasProfile {
		snilsValue = prof.SNILS().Value()
		dobValue = prof.DateOfBirth().Date()
		jobTitleValue = prof.JobTitle().Value()
		educationValue = prof.Education().Value()
		if prof.HasOrganization() {
			organizationID = prof.OrganizationID()
		}
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO persons (
			id, first_name, last_name, middle_name, snils, date_of_birth,
			job_title, education, organization_id, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, ''), COALESCE($8, ''), $9, now())
		ON CONFLICT (id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			middle_name = EXCLUDED.middle_name,
			snils = EXCLUDED.snils,
			date_of_birth = EXCLUDED.date_of_birth,
			job_title = EXCLUDED.job_title,
			education = EXCLUDED.education,
			organization_id = EXCLUDED.organization_id,
			deleted_at = NULL,
			updated_at = now()`,
		p.ID(), p.Name().FirstName(), p.Name().LastName(), p.Name().MiddleName(),
		snilsValue, dobValue, jobTitleValue, educationValue, organizationID)
	if err != nil {
		return fmt.Errorf("save person: %w", err)
	}
	return nil
}

func (r *PersonRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE persons
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}
	return nil
}

func (r *PersonRepository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]personports.PersonShortView, error) {
	return r.list(ctx, `p.organization_id = $1`, organizationID, limit, offset)
}

func (r *PersonRepository) ListByLastName(ctx context.Context, lastName string, limit, offset int) ([]personports.PersonShortView, error) {
	return r.list(ctx, `lower(p.last_name) LIKE lower('%' || $1 || '%')`, lastName, limit, offset)
}

func (r *PersonRepository) ListBySnils(ctx context.Context, snilsValue string, limit, offset int) ([]personports.PersonShortView, error) {
	return r.list(ctx, `p.snils LIKE $1 || '%'`, snilsValue, limit, offset)
}

func (r *PersonRepository) GetDetailsByID(ctx context.Context, id uuid.UUID) (personports.PersonDetailedView, error) {
	var view personports.PersonDetailedView
	var snilsValue, dateOfBirth, organizationName sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT p.id::text, p.first_name, p.last_name, p.middle_name,
			p.snils, to_char(p.date_of_birth, 'YYYY-MM-DD'), p.job_title, p.education,
			o.name
		FROM persons p
		LEFT JOIN organizations o ON o.id = p.organization_id
		WHERE p.id = $1 AND p.deleted_at IS NULL`,
		id).Scan(&view.ID, &view.FirstName, &view.LastName, &view.MiddleName,
		&snilsValue, &dateOfBirth, &view.JobTitle, &view.Education, &organizationName)
	if err != nil {
		return personports.PersonDetailedView{}, err
	}
	view.Snils = snilsValue.String
	view.DateOfBirth = dateOfBirth.String
	view.OrganizationName = organizationName.String
	return view, nil
}

func (r *PersonRepository) RecordPersonAudit(ctx context.Context, event personports.AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal person audit event: %w", err)
	}
	var entityID any
	if event.PersonID != "" {
		entityID = event.PersonID
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO audit_events (action, entity_id, actor_role, payload)
		VALUES ($1, $2, $3, $4)`,
		string(event.Action), entityID, event.ActorRole, payload)
	if err != nil {
		return fmt.Errorf("record person audit: %w", err)
	}
	return nil
}

func (r *PersonRepository) list(ctx context.Context, where string, arg any, limit, offset int) ([]personports.PersonShortView, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT p.id::text,
			trim(p.last_name || ' ' || p.first_name || ' ' || p.middle_name),
			COALESCE(o.name, '')
		FROM persons p
		LEFT JOIN organizations o ON o.id = p.organization_id
		WHERE p.deleted_at IS NULL AND %s
		ORDER BY p.last_name, p.first_name, p.id
		LIMIT $2 OFFSET $3`, where), arg, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []personports.PersonShortView
	for rows.Next() {
		var view personports.PersonShortView
		if err := rows.Scan(&view.ID, &view.FullName, &view.OrganizationName); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}
