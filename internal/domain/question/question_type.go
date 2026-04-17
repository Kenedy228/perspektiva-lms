package question

type Type string

const (
	TypeSelectable Type = "selectable"
	TypeMatching   Type = "matching"
	TypeSequence   Type = "sequence"
	TypeTyped      Type = "typed"
	TypeShort      Type = "short"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeMatching, TypeSelectable, TypeSequence, TypeShort, TypeTyped:
		return true
	default:
		return false
	}
}

func (t Type) String() string {
	return string(t)
}
