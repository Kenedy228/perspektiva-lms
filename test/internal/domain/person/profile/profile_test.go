//go:build legacy
// +build legacy

package profile_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Arrange
	orgID := uuid.New()
	prof, err := profile.New(profile.Params{
		Snils:          snilsFixture(),
		DateOfBirth:    dateOfBirthFixture(),
		JobTitle:       jobTitleFixture(),
		Education:      educationFixture(),
		OrganizationID: orgID,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, snilsFixture(), prof.Snils())
	assert.Equal(t, dateOfBirthFixture(), prof.DateOfBirth())
	assert.Equal(t, jobTitleFixture(), prof.JobTitle())
	assert.Equal(t, educationFixture(), prof.Education())
	assert.Equal(t, orgID, prof.OrganizationID())
}

func TestHasOrganization(t *testing.T) {
	tc := []struct {
		name    string
		orgID   uuid.UUID
		wantHas bool
	}{
		{
			name:    "orgID не равен uuid.Nil",
			orgID:   uuid.New(),
			wantHas: true,
		},
		{
			name:    "orgID равен uuid.Nil",
			orgID:   uuid.Nil,
			wantHas: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			prof, err := profile.New(profile.Params{
				Snils:          snilsFixture(),
				DateOfBirth:    dateOfBirthFixture(),
				JobTitle:       jobTitleFixture(),
				Education:      educationFixture(),
				OrganizationID: tt.orgID,
			})

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.wantHas, prof.HasOrganization())
		})
	}
}

// Fixtures

func snilsFixture() snils.Snils {
	s, _ := snils.New("11223344595")
	return s
}

func dateOfBirthFixture() dob.DateOfBirth {
	db, _ := dob.New(time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)), time.Now())
	return db
}

func jobTitleFixture() jobtitle.JobTitle {
	jt, _ := jobtitle.New("jobtitle fixture")
	return jt
}

func educationFixture() education.Education {
	edu, _ := education.New("education fixture")
	return edu
}
