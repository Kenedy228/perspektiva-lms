package permission

type Resource int

const (
	unknownResource Resource = iota
	ResourceUser
	ResourceCourse
	resourceCount
)

func (r Resource) IsValid() bool {
	return r > unknownResource && r < resourceCount
}
