package permission

import (
	"slices"
)

type PermissionSet struct {
	items map[Resource][]Action
}

func NewSet(permissions []*Permission) (*PermissionSet, error) {
	items := make(map[Resource][]Action, len(permissions))

	for i := range permissions {
		resource := permissions[i].resource
		actions := permissions[i].actions

		if _, ok := items[resource]; ok {
			return nil, ErrPermissionExists
		}

		items[resource] = actions
	}

	return &PermissionSet{
		items: items,
	}, nil
}

func (ps *PermissionSet) Items() []*Permission {
	res := make([]*Permission, 0, len(ps.items))

	for resource, actions := range ps.items {
		permission, _ := New(resource, actions)

		res = append(res, permission)
	}

	return res
}

func (ps *PermissionSet) Has(resource Resource, action Action) (bool, error) {
	if !resource.IsValid() {
		return false, ErrInvalidResource
	}

	if !action.IsValid() {
		return false, ErrInvalidAction
	}

	actions, ok := ps.items[resource]
	if !ok {
		return false, nil
	}

	return slices.Contains(actions, action), nil
}

func (ps *PermissionSet) Add(permission *Permission) error {
	resource := permission.resource
	actions := permission.actions

	if _, ok := ps.items[resource]; ok {
		return ErrPermissionExists
	}

	ps.items[resource] = actions
	return nil
}

func (ps *PermissionSet) Remove(resource Resource) error {
	if !resource.IsValid() {
		return ErrInvalidResource
	}

	if _, ok := ps.items[resource]; !ok {
		return ErrNoResourceFound
	}

	delete(ps.items, resource)
	return nil
}
