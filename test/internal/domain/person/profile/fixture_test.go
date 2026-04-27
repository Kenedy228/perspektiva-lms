package profile_test

import (
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/google/uuid"
)

func snilsFixture() snils.Snils {
	s, _ := snils.New("11223344595")
	return s
}

func dateOfBirthFixture() dob.DateOfBirth {
	db, _ := dob.New(time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)), time.Now())
	return db
}

var orgID = uuid.New()

func organizationIDFixture() uuid.UUID {
	return orgID 
}
