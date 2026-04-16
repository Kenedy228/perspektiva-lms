package resource

type Resource int

const (
	unknown Resource = iota
	ResourceUser
	ResourceCourse
	count
)

func (r Resource) IsValid() bool {
	return r > unknown && r < count
}
