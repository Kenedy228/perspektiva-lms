package permission

import (
	"errors"

	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

var (
	ErrInvalidResource  = errors.New("resource is invalid")
	ErrInvalidAction    = errors.New("action is invalid")
	ErrEmptyActions     = errors.New("empty actions")
	ErrDuplicatedAction = errors.New("duplicated action")
)

func validateResource(resource resource.Resource) error {
	if !resource.IsValid() {
		return ErrInvalidResource
	}

	return nil
}

func validateActions(actions []action.Action) error {
	if len(actions) == 0 {
		return ErrEmptyActions
	}

	visited := make(map[action.Action]struct{}, len(actions))
	for i := range actions {
		a := actions[i]

		if !a.IsValid() {
			return ErrInvalidAction
		}

		if _, ok := visited[a]; ok {
			return ErrDuplicatedAction
		}

		visited[a] = struct{}{}
	}

	return nil
}
