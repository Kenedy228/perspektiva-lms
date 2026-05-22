package question

type Type string

const (
	TypeSelectable Type = "selectable"
	TypeMatching   Type = "matching"
	TypeSequence   Type = "sequence"
	TypeShort      Type = "short"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeSelectable, TypeMatching, TypeSequence, TypeShort:
		return true
	default:
		return false
	}
}

func (t Type) Title() string {
	switch t {
	case TypeMatching:
		return "на соответствие"
	case TypeSelectable:
		return "выбор"
	case TypeSequence:
		return "последовательность"
	case TypeShort:
		return "короткий ответ"
	default:
		return ""
	}
}

func (t Type) DefaultInstruction() string {
	switch t {
	case TypeMatching:
		return "сопоставьте элементы из левого списка с элементами из правого"
	case TypeSelectable:
		return "выберите один или несколько правильных вариантов ответа"
	case TypeSequence:
		return "расставьте события в хронологическом порядке"
	case TypeShort:
		return "введите короткий ответ"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
