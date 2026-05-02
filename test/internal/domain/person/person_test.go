package person_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person"
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"gitflic.ru/lms/internal/domain/person/profile/education"
	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/internal/domain/person/profile/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Arrange
	p, err := person.New(nameFixture())

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, p.ID())
	assert.Equal(t, nameFixture(), p.Name())
	assert.False(t, p.HasProfile())
}

func TestAttachProfile(t *testing.T) {
	// Arrange
	p, err := person.New(nameFixture())
	require.NoError(t, err)

	// Act
	newProf := profileFixture()
	p.AttachProfile(newProf)
	gotProf, ok := p.Profile()

	// Assert
	assert.Equal(t, newProf, gotProf)
	assert.True(t, ok)
	assert.True(t, p.HasProfile())
}

func TestDetachProfile(t *testing.T) {
	t.Run("если профиль установлен, то удаляет профиль", func(t *testing.T) {
		// Arrange
		p, err := person.New(nameFixture())
		require.NoError(t, err)

		// Act
		newProf := profileFixture()
		p.AttachProfile(newProf)

		p.DetachProfile()
		_, ok := p.Profile()

		// Assert
		assert.False(t, p.HasProfile())
		assert.False(t, ok)
	})

	t.Run("если профиль не установлен, то ничего не делает", func(t *testing.T) {
		// Arrange
		p, err := person.New(nameFixture())
		require.NoError(t, err)

		// Act
		p.DetachProfile()
		_, ok := p.Profile()

		// Assert
		assert.False(t, p.HasProfile())
		assert.False(t, ok)
	})
}

func TestRename(t *testing.T) {
	// Arrange
	p, err := person.New(nameFixture())
	require.NoError(t, err)

	// Act
	newName, err := name.New("Петров", "Петр", "Петрович")
	require.NoError(t, err)

	p.Rename(newName)

	// Assert
	assert.Equal(t, newName, p.Name())
}

//Fixtures

func profileFixture() profile.Profile {
	snils, _ := snils.New("11223344595")
	db, _ := dob.New(time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)), time.Now())
	jt, _ := jobtitle.New("ведущий инженер")
	edu, _ := education.New("высшее")

	return profile.New(snils, db, jt, edu, uuid.Nil)
}

func nameFixture() name.Name {
	name, _ := name.New("Иванов", "Иван", "Иванович")
	return name
}
