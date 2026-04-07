package organization

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	id        uuid.UUID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func New(name string) (*Organization, error) {
	if name == "" {
		return nil, ErrEmptyName
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate organization id: %w", err)
	}

	return &Organization{
		id:        id,
		name:      name,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

func (o *Organization) Id() uuid.UUID {
	return o.id
}

func (o *Organization) Name() string {
	return o.name
}

func (o *Organization) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Organization) UpdatedAt() time.Time {
	return o.updatedAt
}

func (o *Organization) Rename(newName string) error {
	if newName == "" {
		return ErrEmptyName
	}

	o.name = newName
	o.updatedAt = time.Now()
	return nil
}

func (o *Organization) Equal(other *Organization) bool {
	if other == nil {
		return false
	}
	return o.id == other.id
}
