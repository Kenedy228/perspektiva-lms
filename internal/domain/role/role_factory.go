package role

import (
	"gitflic.ru/lms/internal/domain/permission"
	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

func getAdminPermissions() []permission.Permission {
	return []permission.Permission{
		makePermission(resource.ResourceCourse, []action.Action{action.ActionRead, action.ActionWrite}),
		makePermission(resource.ResourceUser, []action.Action{action.ActionRead, action.ActionWrite}),
	}
}

func getCreatorPermissions() []permission.Permission {
	return []permission.Permission{
		makePermission(resource.ResourceCourse, []action.Action{action.ActionRead, action.ActionWrite}),
	}
}

func getStudentPermissions() []permission.Permission {
	return []permission.Permission{
		makePermission(resource.ResourceCourse, []action.Action{action.ActionRead}),
	}
}

func getOrganizationPermissions() []permission.Permission {
	return []permission.Permission{
		makePermission(resource.ResourceUser, []action.Action{action.ActionRead}),
	}
}

func makePermission(resource resource.Resource, actions []action.Action) permission.Permission {
	params := permission.Params{
		Resource: resource,
		Actions:  actions,
	}
	p, err := permission.New(params)
	if err != nil {
		panic("developer error: invalid permission configured: " + err.Error())
	}

	return p
}
