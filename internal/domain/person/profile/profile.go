package profile

import (
	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"gitflic.ru/lms/internal/domain/person/profile/education"
	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

type Profile struct {
	snils          snils.Snils
	dateOfBirth    dob.DateOfBirth
	jobTitle       jobtitle.JobTitle
	education      education.Education
	organizationID uuid.UUID
}

func New(s snils.Snils, dateOfBirth dob.DateOfBirth, jt jobtitle.JobTitle, edu education.Education, orgID uuid.UUID) Profile {
	return Profile{
		snils:          s,
		dateOfBirth:    dateOfBirth,
		jobTitle:       jt,
		education:      edu,
		organizationID: orgID,
	}
}

func (p Profile) Snils() snils.Snils {
	return p.snils
}

func (p Profile) DateOfBirth() dob.DateOfBirth {
	return p.dateOfBirth
}

func (p Profile) JobTitle() jobtitle.JobTitle {
	return p.jobTitle
}

func (p Profile) Education() education.Education {
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
