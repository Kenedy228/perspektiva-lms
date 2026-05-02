package profile

import (
	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"gitflic.ru/lms/internal/domain/person/profile/education"
	"gitflic.ru/lms/internal/domain/person/profile/jobtitle"
	"gitflic.ru/lms/internal/domain/person/profile/snils"
	"github.com/google/uuid"
)

type Params struct {
	Snils          snils.Snils
	DateOfBirth    dob.DateOfBirth
	JobTitle       jobtitle.JobTitle
	Education      education.Education
	OrganizationID uuid.UUID
}
