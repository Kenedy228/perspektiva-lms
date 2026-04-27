package profile

import (
	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/google/uuid"
)

type Profile struct {
	snils          snils.Snils
	dateOfBirth    dob.DateOfBirth
	jobTitle       string
	education      string
	organizationID uuid.UUID
}

func New(params Params) (Profile, error) {
	jobTitle := normalize(params.JobTitle)
	education := normalize(params.Education)

	if err := validateJobTitle(jobTitle); err != nil {
		return Profile{}, err
	}

	if err := validateEducation(education); err != nil {
		return Profile{}, err
	}

	return Profile{
		snils:          params.Snils,
		dateOfBirth:    params.DateOfBirth,
		jobTitle:       jobTitle,
		education:      education,
		organizationID: params.OrganizationID,
	}, nil
}

func (p Profile) Snils() snils.Snils {
	return p.snils
}

func (p Profile) DateOfBirth() dob.DateOfBirth {
	return p.dateOfBirth
}

func (p Profile) JobTitle() string {
	return p.jobTitle
}

func (p Profile) Education() string {
	return p.education
}

func (p Profile) OrganizationID() uuid.UUID {
	return p.organizationID
}

func (p Profile) HasOrganization() bool {
	if p.organizationID == uuid.Nil {
		return false
	}

	return true
}

func (p Profile) Clone() Profile {
	return Profile{
		snils:          p.snils,
		dateOfBirth:    p.dateOfBirth,
		jobTitle:       p.jobTitle,
		education:      p.education,
		organizationID: p.organizationID,
	}
}
