package profile

import (
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var orgIDFixture = uuid.New()

func dateOfBirthFixture(t *testing.T) dob.DateOfBirth {
	db, err := dob.New(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
	)

	require.NoError(t, err)
	return db
}

func snilsFixture(t *testing.T) snils.SNILS {
	s, err := snils.New("111-111-111 45")
	require.NoError(t, err)

	return s
}

func educationFixture(t *testing.T) education.Education {
	e, err := education.New("высшее")
	require.NoError(t, err)

	return e
}

func jobTitleFixture(t *testing.T) jobtitle.JobTitle {
	jt, err := jobtitle.New("механик")
	require.NoError(t, err)

	return jt
}
