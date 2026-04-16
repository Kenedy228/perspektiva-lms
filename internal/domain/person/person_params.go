package person

import (
	"gitflic.ru/lms/internal/domain/person/dob"
	"github.com/google/uuid"
)

type Params struct {
	FirstName      string
	LastName       string
	MiddleName     string
	JobTitle       string
	Snils          string
	DateOfBirth    dob.DateOfBirth
	Education      string
	OrganizationID uuid.UUID
}
