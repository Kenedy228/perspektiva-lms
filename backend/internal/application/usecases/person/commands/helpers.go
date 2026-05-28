package commands

import (
	"context"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type clock func() time.Time

func realClock() time.Time {
	return time.Now()
}

type profileInput struct {
	DateOfBirth    time.Time
	Snils          string
	JobTitle       string
	Education      string
	OrganizationID string
}

func buildProfile(in profileInput, now time.Time) (profile.Profile, snils.SNILS, error) {
	db, err := dob.New(in.DateOfBirth, now)
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, fmt.Errorf("create person date of birth: %w", err)
	}

	edu, err := education.New(in.Education)
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, fmt.Errorf("create person education: %w", err)
	}

	jt, err := jobtitle.New(in.JobTitle)
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, fmt.Errorf("create person job title: %w", err)
	}

	s, err := snils.New(in.Snils)
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, fmt.Errorf("create person snils: %w", err)
	}

	orgID, err := parseOptionalUUID(in.OrganizationID, "organization id")
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, err
	}

	prof, err := profile.New(s, db, jt, edu, orgID)
	if err != nil {
		return profile.Profile{}, snils.SNILS{}, fmt.Errorf("create person profile: %w", err)
	}

	return prof, s, nil
}

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
	}

	return id, nil
}

func parseOptionalUUID(value, field string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, nil
	}

	return parseRequiredUUID(value, field)
}

func requireSNILSAvailable(ctx context.Context, r personports.Repository, s snils.SNILS, exclude uuid.UUID) error {
	exists, err := r.SNILSExists(ctx, s, exclude)
	if err != nil {
		return fmt.Errorf("check person snils uniqueness: %w", err)
	}

	if exists {
		return fmt.Errorf("%w: snils already belongs to another person", common.ErrConflict)
	}

	return nil
}

func recordAudit(ctx context.Context, audit personports.AuditRecorder, action personports.AuditAction, personID string, actor role.Role, organizationID string) error {
	if err := audit.RecordPersonAudit(ctx, personports.AuditEvent{
		Action:         action,
		PersonID:       personID,
		ActorRole:      actor.Kind().String(),
		OrganizationID: organizationID,
	}); err != nil {
		return fmt.Errorf("record person audit: %w", err)
	}

	return nil
}
