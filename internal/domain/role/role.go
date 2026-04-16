package role

import (
	"gitflic.ru/lms/internal/domain/permission"
	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

type Role struct {
	roleType    RoleType
	permissions []permission.Permission
}

func NewAsAdmin() Role {
	return Role{
		roleType:    RoleTypeAdmin,
		permissions: getAdminPermissions(),
	}
}

func NewAsCreator() Role {
	return Role{
		roleType:    RoleTypeCreator,
		permissions: getCreatorPermissions(),
	}
}

func NewAsStudent() Role {
	return Role{
		roleType:    RoleTypeStudent,
		permissions: getStudentPermissions(),
	}
}

func NewAsOrganization() Role {
	return Role{
		roleType:    RoleTypeOrganization,
		permissions: getOrganizationPermissions(),
	}
}

func (r Role) Allows(resource resource.Resource, action action.Action) bool {
	for i := range r.permissions {
		if r.permissions[i].Resource() == resource {
			return r.permissions[i].HasAction(action)
		}
	}

	return false
}
