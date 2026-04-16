package action

type Action int

const (
	unknown Action = iota
	ActionRead
	ActionWrite
	count
)

func (a Action) IsValid() bool {
	return a > unknown && a < count
}
