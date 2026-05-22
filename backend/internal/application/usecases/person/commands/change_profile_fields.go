package commands

import (
	"context"
	"fmt"
	"time"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ChangeSNILSUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewChangeSNILSUseCase(r personports.Repository, audit personports.AuditRecorder) *ChangeSNILSUseCase {
	if r == nil {
		panic("person change snils usecase requires repository")
	}
	if audit == nil {
		panic("person change snils usecase requires audit recorder")
	}

	return &ChangeSNILSUseCase{r: r, audit: audit}
}

type ChangeSNILSInput struct {
	ActorRole role.Role
	PersonID  string
	Snils     string
}

type ChangeSNILSOutput struct {
	ID string
}

func (uc *ChangeSNILSUseCase) Execute(ctx context.Context, in ChangeSNILSInput) (*ChangeSNILSOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	s, err := snils.New(in.Snils)
	if err != nil {
		return nil, fmt.Errorf("change person snils: %w", err)
	}

	if err := requireSNILSAvailable(ctx, uc.r, s, pID); err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ChangeSnils(s); err != nil {
		return nil, fmt.Errorf("change person snils aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionChangeSNILS, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &ChangeSNILSOutput{ID: p.ID().String()}, nil
}

type ChangeDateOfBirthUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
	now   clock
}

func NewChangeDateOfBirthUseCase(r personports.Repository, audit personports.AuditRecorder) *ChangeDateOfBirthUseCase {
	if r == nil {
		panic("person change date of birth usecase requires repository")
	}
	if audit == nil {
		panic("person change date of birth usecase requires audit recorder")
	}

	return &ChangeDateOfBirthUseCase{r: r, audit: audit, now: realClock}
}

type ChangeDateOfBirthInput struct {
	ActorRole   role.Role
	PersonID    string
	DateOfBirth time.Time
}

type ChangeDateOfBirthOutput struct {
	ID string
}

func (uc *ChangeDateOfBirthUseCase) Execute(ctx context.Context, in ChangeDateOfBirthInput) (*ChangeDateOfBirthOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	db, err := dob.New(in.DateOfBirth, uc.now())
	if err != nil {
		return nil, fmt.Errorf("change person date of birth: %w", err)
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ChangeDateOfBirth(db); err != nil {
		return nil, fmt.Errorf("change person date of birth aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionChangeDateOfBirth, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &ChangeDateOfBirthOutput{ID: p.ID().String()}, nil
}

type ChangeJobTitleUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewChangeJobTitleUseCase(r personports.Repository, audit personports.AuditRecorder) *ChangeJobTitleUseCase {
	if r == nil {
		panic("person change job title usecase requires repository")
	}
	if audit == nil {
		panic("person change job title usecase requires audit recorder")
	}

	return &ChangeJobTitleUseCase{r: r, audit: audit}
}

type ChangeJobTitleInput struct {
	ActorRole role.Role
	PersonID  string
	JobTitle  string
}

type ChangeJobTitleOutput struct {
	ID string
}

func (uc *ChangeJobTitleUseCase) Execute(ctx context.Context, in ChangeJobTitleInput) (*ChangeJobTitleOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	jt, err := jobtitle.New(in.JobTitle)
	if err != nil {
		return nil, fmt.Errorf("change person job title: %w", err)
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ChangeJobTitle(jt); err != nil {
		return nil, fmt.Errorf("change person job title aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionChangeJobTitle, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &ChangeJobTitleOutput{ID: p.ID().String()}, nil
}

type ChangeEducationUseCase struct {
	r     personports.Repository
	audit personports.AuditRecorder
}

func NewChangeEducationUseCase(r personports.Repository, audit personports.AuditRecorder) *ChangeEducationUseCase {
	if r == nil {
		panic("person change education usecase requires repository")
	}
	if audit == nil {
		panic("person change education usecase requires audit recorder")
	}

	return &ChangeEducationUseCase{r: r, audit: audit}
}

type ChangeEducationInput struct {
	ActorRole role.Role
	PersonID  string
	Education string
}

type ChangeEducationOutput struct {
	ID string
}

func (uc *ChangeEducationUseCase) Execute(ctx context.Context, in ChangeEducationInput) (*ChangeEducationOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	pID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	edu, err := education.New(in.Education)
	if err != nil {
		return nil, fmt.Errorf("change person education: %w", err)
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, fmt.Errorf("find person: %w", err)
	}

	if err := p.ChangeEducation(edu); err != nil {
		return nil, fmt.Errorf("change person education aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("save person: %w", err)
	}

	if err := recordAudit(ctx, uc.audit, personports.AuditActionChangeEducation, p.ID().String(), in.ActorRole, ""); err != nil {
		return nil, err
	}

	return &ChangeEducationOutput{ID: p.ID().String()}, nil
}
