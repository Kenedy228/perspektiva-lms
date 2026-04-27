package organization

import (
	"time"
)

type Organization struct {
	inn       string
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func New(params Params) (*Organization, error) {
	if err := validateName(params.Name); err != nil {
		return nil, err
	}

	if err := validateInn(params.INN); err != nil {
		return nil, err
	}

	now := time.Now()

	return &Organization{
		name:      params.Name,
		inn:       params.INN,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (o *Organization) Name() string {
	return o.name
}

func (o *Organization) INN() string {
	return o.inn
}

func (o *Organization) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Organization) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Organization) Rename(name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	o.name = name
	o.updatedAt = time.Now()
	return nil
}

func (o *Organization) ChangeINN(inn string) error {
	if err := validateInn(inn); err != nil {
		return err
	}

	o.inn = inn
	o.updatedAt = time.Now()
	return nil
}
