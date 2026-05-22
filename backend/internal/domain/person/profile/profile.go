package profile

import (
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

// Profile contains optional personal and employment details for a person.
type Profile struct {
	snils          snils.SNILS
	dateOfBirth    dob.DateOfBirth
	jobTitle       jobtitle.JobTitle
	education      education.Education
	organizationID uuid.UUID
}

// New creates a profile from validated value objects.
func New(
	s snils.SNILS,
	dob dob.DateOfBirth,
	jt jobtitle.JobTitle,
	e education.Education,
	organizationID uuid.UUID,
) (Profile, error) {
	if err := validateSNILS(s); err != nil {
		return Profile{}, err
	}

	if err := validateDateOfBirth(dob); err != nil {
		return Profile{}, err
	}

	return Profile{
		snils:          s,
		dateOfBirth:    dob,
		jobTitle:       jt,
		education:      e,
		organizationID: organizationID,
	}, nil
}

// SNILS returns the person's SNILS value.a
func (p Profile) SNILS() snils.SNILS {
	return p.snils
}

// DateOfBirth returns the person's date of birth.
func (p Profile) DateOfBirth() dob.DateOfBirth {
	return p.dateOfBirth
}

// JobTitle returns the person's optional job title.
func (p Profile) JobTitle() jobtitle.JobTitle {
	return p.jobTitle
}

// Education returns the person's optional education description.
func (p Profile) Education() education.Education {
	return p.education
}

// OrganizationID returns the optional organization identifier.
func (p Profile) OrganizationID() uuid.UUID {
	return p.organizationID
}

// HasOrganization reports whether the profile is attached to an organization.
func (p Profile) HasOrganization() bool {
	if p.organizationID == uuid.Nil {
		return false
	}

	return true
}

// IsZero reports whether the profile has not been initialized.
func (p Profile) IsZero() bool {
	return p.snils.IsZero() || p.dateOfBirth.IsZero()
}
