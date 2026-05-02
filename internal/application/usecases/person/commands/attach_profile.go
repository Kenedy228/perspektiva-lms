package commands

import (
	"context"
	"time"

	personrepo "gitflic.ru/lms/internal/application/ports/person"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"gitflic.ru/lms/internal/domain/person/profile/education"
	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

type AttachProfileUseCase struct {
	r personrepo.Repository
}

func NewAttachProfileUseCase(r personrepo.Repository) *AttachProfileUseCase {
	return &AttachProfileUseCase{
		r: r,
	}
}

type AttachProfileInput struct {
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
	db, err := dob.New(in.DateOfBirth, time.Now())
	if err != nil {
		return nil, err
	}

	edu, err := education.New(in.Education)
	if err != nil {
		return nil, err
	}

	jt, err := jobtitle.New(in.JobTitle)
	if err != nil {
		return nil, err
	}

	s, err := snils.New(in.Snils)
	if err != nil {
		return nil, err
	}

	var orgID uuid.UUID
	if in.OrganizationID == "" {
		orgID = uuid.Nil
	} else {
		orgID, err = uuid.Parse(in.OrganizationID)
		if err != nil {
			return nil, err
		}
	}

	prof := profile.New(s, db, jt, edu, orgID)

	pID, err := uuid.Parse(in.PersonID)
	if err != nil {
		return nil, err
	}

	p, err := uc.r.FindByID(ctx, pID)
	if err != nil {
		return nil, err
	}

	p.AttachProfile(prof)

	err = uc.r.Save(ctx, p)
	if err != nil {
		return nil, err
	}

	return &AttachProfileOutput{
		ID: p.ID().String(),
	}, nil
}
