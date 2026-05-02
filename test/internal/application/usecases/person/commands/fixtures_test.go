package commands_test

import (
	"time"

	"gitflic.ru/lms/internal/domain/person"
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"gitflic.ru/lms/internal/domain/person/profile/education"
	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/internal/domain/person/profile/snils"
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

func jobTitleFixture() jobtitle.JobTitle {
	jt, _ := jobtitle.New("jobtitle fixture")
	return jt
}

func educationFixture() education.Education {
	edu, _ := education.New("education fixture")
	return edu
}

func personFixture() *person.Person {
	n, _ := name.New("Иван", "Иванов", "Иванович")
	p, _ := person.New(n)
	return p
}

func profileFixture() profile.Profile {
	return profile.New(snilsFixture(), dateOfBirthFixture(), jobTitleFixture(), educationFixture(), uuid.Nil)
}

func personWithProfileFixture() *person.Person {
	p := personFixture()
	p.AttachProfile(profileFixture())
	return p
}
