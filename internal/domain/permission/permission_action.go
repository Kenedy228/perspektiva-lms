package permission

type Action int

const (
	unknownAction Action = iota
	ActionRead
	ActionWrite
	actionCount
)

func (a Action) IsValid() bool {
	return a > unknownAction && a < actionCount
}

