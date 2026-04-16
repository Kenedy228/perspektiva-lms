package organization

import (
	"time"

	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Organization struct {
	id        uuid.UUID
	inn       string
	name      string
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func New(params Params) (*Organization, error) {
	if err := validateName(params.Name); err != nil {
		return nil, err
	}

	if err := validateInn(params.Inn); err != nil {
		return nil, err
	}

	id, err := utils.GenerateID()

	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Organization{
		id:        id,
		name:      params.Name,
		inn:       params.Inn,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (o *Organization) Id() uuid.UUID {
	return o.id
}

func (o *Organization) Name() string {
	return o.name
}

func (o *Organization) Inn() string {
	return o.inn
}

func (o *Organization) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Organization) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Organization) DeletedAt() time.Time {
	return o.deletedAt
}

func (o *Organization) Rename(name string) error {
	if err := validateName(name); err != nil {
		return err
	}

	o.name = name
	o.updatedAt = time.Now()
	return nil
}

func (o *Organization) ChangeInn(inn string) error {
	if err := validateInn(inn); err != nil {
		return err
	}

	o.inn = inn
	o.updatedAt = time.Now()
	return nil
}

func (o *Organization) IsDeleted() bool {
	return !o.deletedAt.IsZero()
}

func (o *Organization) Delete() {
	if o.IsDeleted() {
		return
	}

	now := time.Now()
	o.deletedAt = now
	o.updatedAt = now
}

func (o *Organization) Equal(other *Organization) bool {
	if other == nil {
		return false
	}
	return o.id == other.id
}
