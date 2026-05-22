package organization

import (
	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

// Organization represents an external organization that can participate in LMS flows.
type Organization struct {
	id  uuid.UUID
	inn inn.INN
	n   name.Name
}

// New creates an organization with a generated ID.
func New(inn inn.INN, n name.Name) (*Organization, error) {
	if err := validateINN(inn); err != nil {
		return nil, err
	}

	if err := validateName(n); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Organization{
		id:  id,
		inn: inn,
		n:   n,
	}, nil
}

// Restore recreates an existing organization from persisted state.
func Restore(id uuid.UUID, inn inn.INN, n name.Name) (*Organization, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateINN(inn); err != nil {
		return nil, err
	}

	if err := validateName(n); err != nil {
		return nil, err
	}

	return &Organization{
		id:  id,
		inn: inn,
		n:   n,
	}, nil
}

// ID returns the organization identifier.
func (o *Organization) ID() uuid.UUID {
	return o.id
}

// INN returns the organization's tax identifier.
func (o *Organization) INN() inn.INN {
	return o.inn
}

// Name returns the organization name.
func (o *Organization) Name() name.Name {
	return o.n
}

// ChangeINN changes the organization tax identifier.
func (o *Organization) ChangeINN(inn inn.INN) error {
	if err := validateINN(inn); err != nil {
		return err
	}

	o.inn = inn
	return nil
}

// ChangeName changes the organization name.
func (o *Organization) ChangeName(n name.Name) error {
	if err := validateName(n); err != nil {
		return err
	}

	o.n = n
	return nil
}
