package person

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Person struct {
	id             uuid.UUID
	firstName      string
	lastName       string
	middleName     string
	organizationID uuid.UUID
	createdAt      time.Time
	updatedAt      time.Time
}

func New(firstName, lastName, middleName string, organizationID uuid.UUID) (*Person, error) {
	if strings.TrimSpace(firstName) == "" {
		return nil, ErrEmptyFirstName
	}

	if strings.TrimSpace(lastName) == "" {
		return nil, ErrEmptyLastName
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	now := time.Now()

	return &Person{
		id:             id,
		firstName:      firstName,
		lastName:       lastName,
		middleName:     middleName,
		organizationID: organizationID,
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

func (p *Person) HasMiddleName() bool {
	return p.middleName != ""
}

func (p *Person) OrganizationID() uuid.UUID {
	return p.organizationID
}

func (p *Person) HasOrganization() bool {
	return p.organizationID != uuid.Nil
}

func (p *Person) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Person) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Person) Rename(firstName, lastName, middleName string) error {
	if strings.TrimSpace(firstName) == "" {
		return ErrEmptyFirstName
	}

	if strings.TrimSpace(lastName) == "" {
		return ErrEmptyLastName
	}

	p.firstName = firstName
	p.lastName = lastName
	p.middleName = middleName

	p.updatedAt = time.Now()
	return nil
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

func (p *Person) Equal(other *Person) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}
