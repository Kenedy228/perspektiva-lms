package role

import (
	"testing"

	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

func TestAllows(t *testing.T) {
	organization := NewAsOrganization()

	if organization.Allows(resource.ResourceCourse, action.ActionRead) {
		t.Errorf("expected not allows")
	}

	if organization.Allows(resource.ResourceCourse, action.ActionWrite) {
		t.Errorf("expected not allows")
	}

	if !organization.Allows(resource.ResourceUser, action.ActionRead) {
		t.Errorf("expected allows")
	}
}
