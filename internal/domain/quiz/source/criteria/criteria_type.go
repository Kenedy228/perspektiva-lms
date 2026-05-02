package criteria

type Type string

const (
	TypeRandom Type = "random"
	TypeManual Type = "manual"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeRandom, TypeManual:
		return true
	default:
		return false
	}
}

func (t Type) Title() string {
	switch t {
	case TypeRandom:
		return "случайный"
	case TypeManual:
		return "ручной"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
