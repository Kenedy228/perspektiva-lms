package person_test

import (
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/google/uuid"
)

func profileFixture() *profile.Profile {
	snils, _ := snils.New("11223344595")
	db, _ := dob.New(time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)), time.Now())

	params := profile.Params{
		Snils:          snils,
		DateOfBirth:    db,
		OrganizationID: uuid.New(),
		JobTitle:       "ведущий инженер",
		Education:      "высшее",
	}

	profile, _ := profile.New(params)
	return &profile
}

func nameFixture() name.Name {
	params := name.Params{
		Firstname:  "Иванов",
		Lastname:   "Иван",
		Middlename: "Иванович",
	}

	name, _ := name.New(params)
	return name
}
