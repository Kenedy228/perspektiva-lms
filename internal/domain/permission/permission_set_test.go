package permission

import (
	"slices"
	"testing"
)

func TestHas(t *testing.T) {
	type given struct {
		permissions []*Permission
	}

	type when struct {
		resource Resource
		action   Action
	}

	type want struct {
		has bool
		err error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "set has resource and action",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceUser,
				action:   ActionRead,
			},
			want: want{
				has: true,
				err: nil,
			},
		},
		{
			name: "set has not resource",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceCourse,
				action:   ActionRead,
			},
			want: want{
				has: false,
				err: nil,
			},
		},
		{
			name: "set has resource, but has not action",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceUser,
				action:   ActionWrite,
			},
			want: want{
				has: false,
				err: nil,
			},
		},
		{
			name: "resource invalid",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: Resource(1000),
				action:   ActionWrite,
			},
			want: want{
				has: false,
				err: ErrInvalidResource,
			},
		},
		{
			name: "action invalid",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceUser,
				action:   Action(5000),
			},
			want: want{
				has: false,
				err: ErrInvalidAction,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, _ := NewSet(tt.given.permissions)

			res, err := ps.Has(tt.when.resource, tt.when.action)

			if err != tt.want.err {
				t.Fatalf("want err %v, got %v", tt.want.err, err)
			}

			if res != tt.want.has {
				t.Fatalf("want res %v, got %v", tt.want.has, res)
			}
		})
	}
}

func TestItems(t *testing.T) {
	type given struct {
		permissions []*Permission
	}

	type want struct {
		permissions []*Permission
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "set has resource and action",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
					{
						resource: ResourceCourse,
						actions:  []Action{ActionRead, ActionWrite},
					},
				},
			},
			want: want{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
					{
						resource: ResourceCourse,
						actions:  []Action{ActionRead, ActionWrite},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, _ := NewSet(tt.given.permissions)

			got := ps.Items()

			if len(got) != len(tt.want.permissions) {
				t.Fatalf("expected length of Items() equal to %d, got %d", len(tt.want.permissions), len(got))
			}

			for i := range got {
				var found bool

				for j := range tt.want.permissions {
					if got[i].resource != tt.want.permissions[j].resource {
						continue
					}

					if !slices.Equal(got[i].actions, tt.want.permissions[j].actions) {
						continue
					}

					found = true
				}

				if !found {
					t.Fatalf("can't find corresponding permission in Items and want")
				}
			}
		})
	}
}

func TestAdd(t *testing.T) {
	type given struct {
		permissions []*Permission
	}

	type when struct {
		permission *Permission
	}

	type want struct {
		len int
		err error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "resource doesn't exist",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				permission: &Permission{
					resource: ResourceCourse,
					actions:  []Action{ActionRead},
				},
			},
			want: want{
				len: 2,
				err: nil,
			},
		},
		{
			name: "resource exists",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				permission: &Permission{
					resource: ResourceUser,
					actions:  []Action{ActionRead},
				},
			},
			want: want{
				len: 1,
				err: ErrPermissionExists,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, _ := NewSet(tt.given.permissions)

			if err := ps.Add(tt.when.permission); err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if tt.want.len != len(ps.items) {
				t.Fatalf("expected len of items %d, got %d", tt.want.len, len(ps.items))
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type given struct {
		permissions []*Permission
	}

	type when struct {
		resource Resource
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "resource exists",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceUser,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "resource doesn't exists",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: ResourceCourse,
			},
			want: want{
				err: ErrNoResourceFound,
			},
		},
		{
			name: "resource invalid",
			given: given{
				permissions: []*Permission{
					{
						resource: ResourceUser,
						actions:  []Action{ActionRead},
					},
				},
			},
			when: when{
				resource: Resource(1000),
			},
			want: want{
				err: ErrInvalidResource,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps, _ := NewSet(tt.given.permissions)

			if err := ps.Remove(tt.when.resource); err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if _, ok := ps.items[tt.when.resource]; ok {
				t.Fatalf("expected delete item by %v", tt.when.resource)
			}
		})
	}
}
