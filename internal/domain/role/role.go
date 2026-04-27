package role

type Role struct {
	kind Type
}

func NewAdmin() Role {
	return Role{
		kind: TypeAdmin,
	}
}

func NewCreator() Role {
	return Role{
		kind: TypeCreator,
	}
}

func NewStudent() Role {
	return Role{
		kind: TypeStudent,
	}
}

func NewOrganization() Role {
	return Role{
		kind: TypeOrganization,
	}
}

func (r Role) Kind() Type {
	return r.kind
}

func (r Role) Allows(resource Resource, action Action) bool {
	return r.kind.Allows(resource, action)
}
