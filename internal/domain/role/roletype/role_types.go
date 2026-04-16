package roletype

type RoleType int

const (
	unknown RoleType = iota
	TypeAdmin
	TypeCreator
	TypeStudent
	TypeOrganization
	count
)

func (t RoleType) IsValid() bool {
	return t > unknown && t < count
}
