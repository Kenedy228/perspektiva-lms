package role

import "slices"

type Type string

const (
	TypeAdmin        Type = "admin"
	TypeCreator      Type = "creator"
	TypeStudent      Type = "student"
	TypeOrganization Type = "organization"
)

var (
	permissionsByType = map[Type]map[Resource][]Action{
		TypeAdmin: {
			ResourceCourse: {ActionRead, ActionWrite},
			ResourceUser:   {ActionRead, ActionWrite},
		},
		TypeCreator: {
			ResourceCourse: {ActionRead, ActionWrite},
		},
		TypeStudent: {
			ResourceCourse: {ActionRead},
		},
		TypeOrganization: {
			ResourceUser: {ActionRead},
		},
	}
)

func (t Type) Allows(resource Resource, action Action) bool {
	resources, ok := permissionsByType[t]
	if !ok {
		return false
	}

	actions, ok := resources[resource]
	if !ok {
		return false
	}

	return slices.Contains(actions, action)
}

func (t Type) Title() string {
	switch t {
	case TypeAdmin:
		return "администратор"
	case TypeCreator:
		return "создатель"
	case TypeStudent:
		return "студент"
	case TypeOrganization:
		return "организация"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
