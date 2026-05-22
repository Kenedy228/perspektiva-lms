package commands

import (
	"context"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type AttachProfileUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
	now   clock
}

func NewAttachProfileUseCase(r personports.Repository, audit personports.AuditRecorder) *AttachProfileUseCase {
	if r == nil {
		panic("person attach profile usecase requires repository")
	}
	if audit == nil {
		panic("person attach profile usecase requires audit recorder")
	}

	return &AttachProfileUseCase{
		r:     r,
		audit: audit,
		now:   realClock,
	}
}

type AttachProfileInput struct {
	ActorRole      role.Role
	DateOfBirth    time.Time
	PersonID       string
	Snils          string
	JobTitle       string
	Education      string
	OrganizationID string
}

type AttachProfileOutput struct {
	ID string
}

func (uc *AttachProfileUseCase) Execute(ctx context.Context, in AttachProfileInput) (*AttachProfileOutput, error) {
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

	if err := p.AttachOrReplaceProfile(prof); err != nil {
		return nil, fmt.Errorf("attach person profile: %w", err)
	}

	err = uc.r.Save(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionAttachProfile, p.ID().String(), in.ActorRole, in.OrganizationID); err != nil {
		return nil, err
	}

	return &AttachProfileOutput{
		ID: p.ID().String(),
	}, nil
}
