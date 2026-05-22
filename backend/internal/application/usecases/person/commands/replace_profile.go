package commands

import (
	"context"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ReplaceProfileUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
	now   clock
}

func NewReplaceProfileUseCase(r personports.Repository, audit personports.AuditRecorder) *ReplaceProfileUseCase {
	if r == nil {
		panic("person replace profile usecase requires repository")
	}
	if audit == nil {
		panic("person replace profile usecase requires audit recorder")
	}

	return &ReplaceProfileUseCase{
		r:     r,
		audit: audit,
		now:   realClock,
	}
}

type ReplaceProfileInput struct {
	ActorRole      role.Role
	DateOfBirth    time.Time
	PersonID       string
	Snils          string
	JobTitle       string
	Education      string
	OrganizationID string
}

type ReplaceProfileOutput struct {
	ID string
}

func (uc *ReplaceProfileUseCase) Execute(ctx context.Context, in ReplaceProfileInput) (*ReplaceProfileOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	prof, s, err := buildProfile(profileInput{
		DateOfBirth:    in.DateOfBirth,
		Snils:          in.Snils,
		JobTitle:       in.JobTitle,
		Education:      in.Education,
		OrganizationID: in.OrganizationID,
	}, uc.now())
	if err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	if err := requireSNILSAvailable(ctx, uc.r, s, pID); err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ReplaceProfile(prof); err != nil {
		return nil, fmt.Errorf("replace person profile: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionReplaceProfile, p.ID().String(), in.ActorRole, in.OrganizationID); err != nil {
		return nil, err
	}

	return &ReplaceProfileOutput{ID: p.ID().String()}, nil
}
