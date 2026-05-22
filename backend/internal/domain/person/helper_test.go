package person

import (
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var idFixture = uuid.New()

func nameFixture(t *testing.T) name.Name {
	n, err := name.New("Иван", "Иванов", "Иванович")
	require.NoError(t, err)
	return n
}

func profileFixture(t *testing.T) profile.Profile {
	db, err := dob.New(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
	)
	require.NoError(t, err)

	s, err := snils.New("111-111-111 45")
	require.NoError(t, err)

	prof, err := profile.New(
		s,
		db,
		jobtitle.JobTitle{},
		education.Education{},
		uuid.Nil,
	)
	require.NoError(t, err)

	return prof
}
