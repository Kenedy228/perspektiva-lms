package role

import (
	"errors"
	"testing"
	"unsafe"

	"gitflic.ru/lms/internal/domain/permission"
)

func TestNew(t *testing.T) {
	type given struct {
		name        string
		permissions []*permission.Permission
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "empty name error",
			given: given{
				name:        "",
				permissions: []*permission.Permission{},
			},
			want: want{
				err: ErrEmptyName,
			},
		},
		{
			name: "should create",
			given: given{
				name:        "admin",
				permissions: []*permission.Permission{},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := New(tt.given.name, tt.given.permissions); !errors.Is(err, tt.want.err) {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestRename(t *testing.T) {
	type given struct {
		name        string
		permissions []*permission.Permission
	}

	type when struct {
		name string
	}

	type want struct {
		name string
		err  error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "should rename",
			given: given{
				name:        "old name",
				permissions: []*permission.Permission{},
			},
			when: when{
				name: "new name",
			},
			want: want{
				name: "new name",
				err:  nil,
			},
		},
		{
			name: "empty name",
			given: given{
				name:        "old name",
				permissions: []*permission.Permission{},
			},
			when: when{
				name: "",
			},
			want: want{
				name: "old name",
				err:  ErrEmptyName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, _ := New(tt.given.name, tt.given.permissions)

			err := role.Rename(tt.when.name)

			if err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if role.Name() != tt.want.name {
				t.Fatalf("expected name %v, got %v", tt.want.name, role.Name())
			}
		})
	}
}

func TestChangePermissions(t *testing.T) {
	type given struct {
		name     string
		resource permission.Resource
		actions  []permission.Action
	}

	type when struct {
		resource permission.Resource
		actions  []permission.Action
	}

	type want struct {
		equality bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "should change",
			given: given{
				name:     "admin",
				resource: permission.ResourceCourse,
				actions:  []permission.Action{permission.ActionRead},
			},
			when: when{
				resource: permission.ResourceCourse,
				actions:  []permission.Action{permission.ActionRead},
			},
			want: want{
				equality: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenPermission, _ := permission.New(tt.given.resource, tt.given.actions)
			whenPermission, _ := permission.New(tt.when.resource, tt.when.actions)

			role, _ := New(tt.given.name, []*permission.Permission{givenPermission})

			ptr1 := unsafe.SliceData(role.Permissions())

			if err := role.ChangePermissions([]*permission.Permission{whenPermission}); err != nil {
				t.Errorf("expected no errors, got %v", err)
			}

			ptr2 := unsafe.SliceData(role.Permissions())

			if (ptr1 == ptr2) != tt.want.equality {
				t.Errorf("expected equality of permissions %v, got %v", tt.want.equality, ptr1 == ptr2)
			}
		})
	}
}

func TestAllows(t *testing.T) {
	type given struct {
		name     string
		resource permission.Resource
		actions  []permission.Action
	}

	type when struct {
		resource permission.Resource
		action   permission.Action
	}

	type want struct {
		has bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "contains action",
			given: given{
				name:     "admin",
				resource: permission.ResourceCourse,
				actions:  []permission.Action{permission.ActionRead},
			},
			when: when{
				resource: permission.ResourceCourse,
				action:   permission.ActionRead,
			},
			want: want{
				has: true,
			},
		},
		{
			name: "doesn't contain action",
			given: given{
				name:     "admin",
				resource: permission.ResourceCourse,
				actions:  []permission.Action{permission.ActionRead},
			},
			when: when{
				resource: permission.ResourceCourse,
				action:   permission.ActionWrite,
			},
			want: want{
				has: false,
			},
		},
		{
			name: "doesn't contain resource",
			given: given{
				name:     "admin",
				resource: permission.ResourceCourse,
				actions:  []permission.Action{permission.ActionRead},
			},
			when: when{
				resource: permission.ResourceUser,
				action:   permission.ActionRead,
			},
			want: want{
				has: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenPermission, _ := permission.New(tt.given.resource, tt.given.actions)

			role, _ := New(tt.given.name, []*permission.Permission{givenPermission})

			if has := role.Allows(tt.when.resource, tt.when.action); has != tt.want.has {
				t.Fatalf("expected allows %v, got %v", tt.want.has, has)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type given struct {
		name        string
		permissions []*permission.Permission
	}

	type when struct {
		name        string
		permissions []*permission.Permission
	}

	type want struct {
		equality bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "not equal",
			given: given{
				name:        "first",
				permissions: []*permission.Permission{},
			},
			when: when{
				name:        "second",
				permissions: []*permission.Permission{},
			},
			want: want{
				equality: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenRole, _ := New(tt.given.name, tt.given.permissions)
			whenRole, _ := New(tt.when.name, tt.when.permissions)

			if equality := givenRole.Equal(whenRole); equality != tt.want.equality {
				t.Fatalf("expected equality %v, got %v", tt.want.equality, equality)
			}
		})
	}
}
