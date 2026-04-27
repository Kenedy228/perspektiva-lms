package profile_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/person/profile"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("with organization", func(t *testing.T) {
		//Arrange
		p := newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withOrganizationID().
			withJobTitle("  ведущий   инженер  ").
			withEducation("  высшее   образование ").
			build(t, nil)

		//Assert
		assert.Equal(t, snilsFixture(), p.Snils())
		assert.Equal(t, dateOfBirthFixture(), p.DateOfBirth())
		assert.Equal(t, "ведущий инженер", p.JobTitle())
		assert.Equal(t, "высшее образование", p.Education())
		assert.Equal(t, p.OrganizationID(), organizationIDFixture())
		assert.True(t, p.HasOrganization())
	})

	t.Run("without organization", func(t *testing.T) {
		//Arrange
		p := newProfileBuilder().
			withSnils().
			withDateOfBirth().
			// organizationID не задаём — остаётся uuid.Nil
			withJobTitle("инженер").
			withEducation("среднее профессиональное").
			build(t, nil)

		//Assert
		assert.Equal(t, snilsFixture(), p.Snils())
		assert.Equal(t, dateOfBirthFixture(), p.DateOfBirth())
		assert.Equal(t, "инженер", p.JobTitle())
		assert.Equal(t, "среднее профессиональное", p.Education())
		assert.Equal(t, uuid.Nil, p.OrganizationID())
		assert.False(t, p.HasOrganization())
	})
}

func TestNewProfileValidationErrors(t *testing.T) {
	t.Run("empty job title after normalize", func(t *testing.T) {
		//Arrange-Assert
		newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withJobTitle("   \t \n ").
			withEducation("высшее").
			build(t, profile.ErrInvalid)
	})

	t.Run("empty education after normalize", func(t *testing.T) {
		//Arrange-Assert
		newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withJobTitle("инженер").
			withEducation("   ").
			build(t, profile.ErrInvalid)
	})

	t.Run("job title too long", func(t *testing.T) {
		//Arrange-Assert
		longTitle := strings.Repeat("а", 1e5+1)

		newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withJobTitle(longTitle).
			withEducation("высшее").
			build(t, profile.ErrInvalid)
	})

	t.Run("education too long", func(t *testing.T) {
		//Arrange-Assert
		longEdu := strings.Repeat("б", 1e5+1)

		newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withJobTitle("инженер").
			withEducation(longEdu).
			build(t, profile.ErrInvalid)
	})
}

func TestClone(t *testing.T) {
	t.Run("testing clone values", func(t *testing.T) {
		//Arrange
		original := newProfileBuilder().
			withSnils().
			withDateOfBirth().
			withJobTitle("инженер").
			withEducation("высшее").
			build(t, nil)

		//Act
		cloned := original.Clone()

		//Assert
		assert.Equal(t, original.Snils(), cloned.Snils())
		assert.Equal(t, original.DateOfBirth(), cloned.DateOfBirth())
		assert.Equal(t, original.JobTitle(), cloned.JobTitle())
		assert.Equal(t, original.Education(), cloned.Education())
		assert.Equal(t, original.OrganizationID(), cloned.OrganizationID())
	})
}
