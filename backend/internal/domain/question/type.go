package question

import (
	"fmt"
	"strings"
)

type Type string

const (
	TypeSelectable Type = "selectable"
	TypeMatching   Type = "matching"
	TypeSequence   Type = "sequence"
	TypeTyped      Type = "typed"
	TypeShort      Type = "short"
)

var validTypes = map[Type]struct{}{
	TypeSelectable: {},
	TypeMatching:   {},
	TypeSequence:   {},
	TypeTyped:      {},
	TypeShort:      {},
}

func ParseType(value string) (Type, error) {
	t := Type(strings.TrimSpace(value))
	if !t.IsValid() {
		return "", fmt.Errorf("unknown question type %q", value)
	}

	return t, nil
}

func (t Type) IsValid() bool {
	_, ok := validTypes[t]
	return ok
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
	case TypeTyped:
		return "пропуски"
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
	case TypeTyped:
		return "заполните пропуски"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
