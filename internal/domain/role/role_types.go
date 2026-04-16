package role

type RoleType int

const (
	unknown RoleType = iota
	RoleTypeAdmin
	RoleTypeCreator
	RoleTypeStudent
	RoleTypeOrganization
	count
)

func (t RoleType) IsValid() bool {
	return t > unknown && t < count
}
