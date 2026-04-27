package profile

import (
	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/google/uuid"
)

type Params struct {
	Snils          snils.Snils
	DateOfBirth    dob.DateOfBirth
	JobTitle       string
	Education      string
	OrganizationID uuid.UUID
}
