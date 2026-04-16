package role

import (
	"testing"

	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
	"github.com/stretchr/testify/assert"
)

func TestAllows(t *testing.T) {
	organization := NewAsOrganization()

	assert.False(t, organization.Allows(resource.ResourceCourse, action.ActionRead), "expected organization to not allow read on course")
	assert.False(t, organization.Allows(resource.ResourceCourse, action.ActionWrite), "expected organization to not allow write on course")

	assert.True(t, organization.Allows(resource.ResourceUser, action.ActionRead), "expected organization to allow read on user")
}

