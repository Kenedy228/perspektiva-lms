package person

import (
	"fmt"
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Person struct {
	id             uuid.UUID
	firstName      string
	lastName       string
	middleName     string
	jobTitle       string
	snils          string
	dateOfBirth    dob.DateOfBirth
	education      string
	organizationID uuid.UUID
	createdAt      time.Time
	updatedAt      time.Time
	deletedAt      time.Time
}

func New(params Params) (*Person, error) {
	if err := validateFirstName(params.FirstName); err != nil {
		return nil, err
	}

	if err := validateLastName(params.LastName); err != nil {
		return nil, err
	}

	if err := validateMiddleName(params.MiddleName); err != nil {
		return nil, err
	}

	if err := validateJobTitle(params.JobTitle); err != nil {
		return nil, err
	}

	if err := validateSnils(params.Snils); err != nil {
		return nil, err
	}

	if err := validateEducation(params.Education); err != nil {
		return nil, err
	}

	id, err := utils.GenerateID()

	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Person{
		id:             id,
		firstName:      params.FirstName,
		lastName:       params.LastName,
		middleName:     params.MiddleName,
		jobTitle:       params.JobTitle,
		education:      params.Education,
		snils:          params.Snils,
		organizationID: params.OrganizationID,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func (p *Person) ID() uuid.UUID {
	return p.id
}

func (p *Person) FirstName() string {
	return p.firstName
}

func (p *Person) LastName() string {
	return p.lastName
}

func (p *Person) MiddleName() string {
	return p.middleName
}

func (p *Person) JobTitle() string {
	return p.jobTitle
}

func (p *Person) Snils() string {
	return p.snils
}

func (p *Person) DateOfBirth() dob.DateOfBirth {
	return p.dateOfBirth
}

func (p *Person) Education() string {
	return p.education
}

func (p *Person) OrganizationID() uuid.UUID {
	return p.organizationID
}

func (p *Person) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Person) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Person) DeletedAt() time.Time {
	return p.deletedAt
}

func (p *Person) Rename(firstName, lastName, middleName string) error {
	if err := validateFirstName(firstName); err != nil {
		return err
	}

	if err := validateLastName(lastName); err != nil {
		return err
	}

	if err := validateMiddleName(middleName); err != nil {
		return err
	}

	p.firstName = firstName
	p.lastName = lastName
	p.middleName = middleName

	p.updatedAt = time.Now()
	return nil
}

func (p *Person) ChangeJobTitle(jobTitle string) error {
	if err := validateJobTitle(jobTitle); err != nil {
		return err
	}

	p.jobTitle = jobTitle
	p.updatedAt = time.Now()
	return nil
}

func (p *Person) ChangeEducation(education string) error {
	if err := validateEducation(education); err != nil {
		return err
	}

	p.education = education
	p.updatedAt = time.Now()
	return nil
}

func (p *Person) ChangeSnils(snils string) error {
	if err := validateSnils(snils); err != nil {
		return err
	}

	p.snils = snils
	p.updatedAt = time.Now()
	return nil
}

func (p *Person) ChangeDateOfBirth(dob dob.DateOfBirth) {
	p.dateOfBirth = dob
	p.updatedAt = time.Now()
}

func (p *Person) ChangeOrganization(organizationID uuid.UUID) {
	if organizationID == uuid.Nil {
		p.RemoveOrganization()
		return
	}
	p.organizationID = organizationID
	p.updatedAt = time.Now()
}

func (p *Person) RemoveOrganization() {
	p.organizationID = uuid.Nil
	p.updatedAt = time.Now()
}

func (p *Person) HasOrganization() bool {
	return p.organizationID != uuid.Nil
}

func (p *Person) Equal(other *Person) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}

func (p *Person) IsDeleted() bool {
	return !p.deletedAt.IsZero()
}

func (p *Person) Delete() {
	if p.IsDeleted() {
		return
	}

	now := time.Now()
	p.updatedAt = now
	p.deletedAt = now
}

func (p *Person) FullName() string {
	if p.middleName == "" {
		return fmt.Sprintf("%s %s", p.firstName, p.lastName)
	}
	return fmt.Sprintf("%s %s %s", p.firstName, p.middleName, p.lastName)
}
