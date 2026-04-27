package profile_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type profileBuilder struct {
	snils          snils.Snils
	dateOfBirth    dob.DateOfBirth
	organizationID uuid.UUID
	jobTitle       string
	education      string
}

func newProfileBuilder() *profileBuilder {
	return &profileBuilder{}
}

func (b *profileBuilder) withSnils() *profileBuilder {
	b.snils = snilsFixture()
	return b
}

func (b *profileBuilder) withDateOfBirth() *profileBuilder {
	b.dateOfBirth = dateOfBirthFixture()
	return b
}

func (b *profileBuilder) withOrganizationID() *profileBuilder {
	b.organizationID = organizationIDFixture()
	return b
}

func (b *profileBuilder) withJobTitle(title string) *profileBuilder {
	b.jobTitle = title
	return b
}

func (b *profileBuilder) withEducation(education string) *profileBuilder {
	b.education = education
	return b
}

func (b *profileBuilder) build(t *testing.T, wantErr error) profile.Profile {
	params := profile.Params{
		DateOfBirth:    b.dateOfBirth,
		OrganizationID: b.organizationID,
		JobTitle:       b.jobTitle,
		Education:      b.education,
		Snils:          b.snils,
	}

	pr, err := profile.New(params)
	assert.ErrorIs(t, err, wantErr)

	return pr
}
