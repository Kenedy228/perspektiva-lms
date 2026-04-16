package permission

import (
	"slices"

	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

type Permission struct {
	resource resource.Resource
	actions  []action.Action
}

func New(params Params) (Permission, error) {
	if err := validateResource(params.Resource); err != nil {
		return Permission{}, err
	}

	if err := validateActions(params.Actions); err != nil {
		return Permission{}, err
	}

	cActions := slices.Clone(params.Actions)

	return Permission{
		resource: params.Resource,
		actions:  cActions,
	}, nil
}

func (p Permission) Resource() resource.Resource {
	return p.resource
}

func (p Permission) Actions() []action.Action {
	return slices.Clone(p.actions)
}

func (p Permission) HasAction(action action.Action) bool {
	return slices.Contains(p.actions, action)
}
