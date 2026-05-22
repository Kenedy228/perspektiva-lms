package commands_test

import (
	"time"

	persondomain "gitflic.ru/lms/backend/internal/domain/person"
	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/person/profile/dob"
	"gitflic.ru/lms/backend/internal/domain/person/profile/education"
	"gitflic.ru/lms/backend/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

func validDOB() time.Time {
	return time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC)
}

func validSNILS() string {
	return "11223344595"
}

func personFixture() *persondomain.Person {
	n, err := name.New("Иван", "Иванов", "Иванович")
	if err != nil {
		panic(err)
	}

	p, err := persondomain.New(n)
	if err != nil {
		panic(err)
	}

	return p
}

func profileFixture(orgID uuid.UUID) profile.Profile {
	s, err := snils.New(validSNILS())
	if err != nil {
		panic(err)
	}

	db, err := dob.New(validDOB(), time.Now())
	if err != nil {
		panic(err)
	}

	jt, err := jobtitle.New("ведущий инженер")
	if err != nil {
		panic(err)
	}

	edu, err := education.New("высшее")
	if err != nil {
		panic(err)
	}

	prof, err := profile.New(profile.Params{
		Snils:          s,
		DateOfBirth:    db,
		JobTitle:       jt,
		Education:      edu,
		OrganizationID: orgID,
	})
	if err != nil {
		panic(err)
	}

	return prof
}

func personWithProfileFixture(orgID uuid.UUID) *persondomain.Person {
	p := personFixture()
	if err := p.AttachOrReplaceProfile(profileFixture(orgID)); err != nil {
		panic(err)
	}

	return p
}
