package role

import (
	"fmt"
	"strings"
	"time"

	"gitflic.ru/lms/internal/domain/permission"
	"github.com/google/uuid"
)

type Role struct {
	id          uuid.UUID
	name        string
	permissions *permission.PermissionSet
	createdAt   time.Time
	updatedAt   time.Time
}

func New(name string, permissions []*permission.Permission) (*Role, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	set, err := permission.NewSet(permissions)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	now := time.Now()

	return &Role{
		id:          id,
		name:        name,
		permissions: set,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Permissions() []*permission.Permission {
	return r.permissions.Items()
}

func (r *Role) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}

	r.name = name
	r.updatedAt = time.Now()
	return nil
}

func (r *Role) ChangePermissions(permissions []*permission.Permission) error {
	set, err := permission.NewSet(permissions)
	if err != nil {
		return err
	}

	r.permissions = set
	r.updatedAt = time.Now()
	return nil
}

func (r *Role) Allows(resource permission.Resource, action permission.Action) bool {
	allowed, err := r.permissions.Has(resource, action)
	if err != nil {
		return false
	}

	return allowed
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Role) Equal(other *Role) bool {
	if other == nil {
		return false
	}

	return r.id == other.id
}
