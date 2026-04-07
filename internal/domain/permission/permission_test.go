package permission

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	type when struct {
		resource Resource
		actions  []Action
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		when
		want
	}{
		{
			name: "valid resource and actions",
			when: when{
				resource: ResourceUser,
				actions:  []Action{ActionRead},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "invalid resource",
			when: when{
				resource: Resource(1000),
				actions:  []Action{ActionRead},
			},
			want: want{
				err: ErrInvalidResource,
			},
		},
		{
			name: "invalid action",
			when: when{
				resource: ResourceUser,
				actions:  []Action{Action(1000)},
			},
			want: want{
				err: ErrInvalidAction,
			},
		},
		{
			name: "duplicated actions",
			when: when{
				resource: ResourceUser,
				actions:  []Action{ActionRead, ActionRead},
			},
			want: want{
				err: ErrDuplicatedActions,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.resource, tt.actions)

			if err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestActions(t *testing.T) {
	type given struct {
		resource Resource
		actions  []Action
	}

	type want struct {
		actions []Action
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "should not change actions",
			given: given{
				resource: ResourceUser,
				actions:  []Action{ActionRead},
			},
			want: want{
				actions: []Action{ActionRead},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permission, _ := New(tt.given.resource, tt.given.actions)

			got := permission.Actions()
			got[0] = ActionWrite

			if pActions := permission.Actions(); !slices.Equal(pActions, tt.want.actions) {
				t.Fatalf("expected slices equal, got false")
			}
		})
	}
}
