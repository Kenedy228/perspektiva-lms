package permission

import (
	"testing"

	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		resource resource.Resource
		actions  []action.Action
		err      error
	}{
		{
			name:     "valid resource and actions",
			resource: resource.ResourceUser,
			actions:  []action.Action{action.ActionRead},
			err:      nil,
		},
		{
			name:     "invalid resource",
			resource: resource.Resource(1000),
			actions:  []action.Action{action.ActionRead},
			err:      ErrInvalidResource,
		},
		{
			name:     "invalid action",
			resource: resource.ResourceUser,
			actions:  []action.Action{action.Action(1000)},
			err:      ErrInvalidAction,
		},
		{
			name:     "duplicated action",
			resource: resource.ResourceUser,
			actions:  []action.Action{action.ActionRead, action.ActionRead},
			err:      ErrDuplicatedAction,
		},
		{
			name:     "empty actions",
			resource: resource.ResourceUser,
			actions:  []action.Action{},
			err:      ErrEmptyActions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(tt.resource, tt.actions)
			p, err := New(params)

			if err != tt.err {
				t.Fatalf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if &p.Actions()[0] == &params.Actions[0] {
					t.Errorf("expected different slices, got the same")
				}
			}
		})
	}
}

func createParams(resource resource.Resource, actions []action.Action) Params {
	return Params{
		Resource: resource,
		Actions:  actions,
	}
}
