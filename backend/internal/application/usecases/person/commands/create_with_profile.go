package commands

import (
	"context"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	persondomain "gitflic.ru/lms/backend/internal/domain/person"
	personname "gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type CreateWithProfileUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
	now   clock
}

func NewCreateWithProfileUseCase(r personports.Repository, audit personports.AuditRecorder) *CreateWithProfileUseCase {
	if r == nil {
		panic("person create with profile usecase requires repository")
	}
	if audit == nil {
		panic("person create with profile usecase requires audit recorder")
	}

	return &CreateWithProfileUseCase{
		r:     r,
		audit: audit,
		now:   realClock,
	}
}

type CreateWithProfileInput struct {
	ActorRole      role.Role
	DateOfBirth    time.Time
	FirstName      string
	LastName       string
	MiddleName     string
	Snils          string
	JobTitle       string
	Education      string
	OrganizationID string
}

type CreateWithProfileOutput struct {
	ID string
}

func (uc *CreateWithProfileUseCase) Execute(ctx context.Context, in CreateWithProfileInput) (*CreateWithProfileOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pName, err := personname.New(in.FirstName, in.LastName, in.MiddleName)
	if err != nil {
		return nil, fmt.Errorf("create person name: %w", err)
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

	if err := requireSNILSAvailable(ctx, uc.r, s, uuid.Nil); err != nil {
		return nil, err
	}

	p, err := persondomain.New(pName)
	if err != nil {
		return nil, fmt.Errorf("create person aggregate: %w", err)
	}

	if err := p.AttachOrReplaceProfile(prof); err != nil {
		return nil, fmt.Errorf("attach person profile: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionCreateWithProfile, p.ID().String(), in.ActorRole, in.OrganizationID); err != nil {
		return nil, err
	}

	return &CreateWithProfileOutput{
		ID: p.ID().String(),
	}, nil
}
