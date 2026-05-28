package role

import (
	"fmt"
	"strings"
)

type Type string

const (
	TypeAdmin        Type = "admin"
	TypeCreator      Type = "creator"
	TypeStudent      Type = "student"
	TypeOrganization Type = "organization"
)

var (
	validTypes = map[Type]struct{}{
		TypeAdmin:        {},
		TypeCreator:      {},
		TypeStudent:      {},
		TypeOrganization: {},
	}
)

func ParseType(value string) (Type, error) {
	t := Type(strings.TrimSpace(value))
	if !t.IsValid() {
		return "", fmt.Errorf("%w: неизвестный тип роли %q", ErrInvalid, value)
	}

	return t, nil
}

func (t Type) IsValid() bool {
	_, ok := validTypes[t]
	return ok
}

func (t Type) Title() string {
	switch t {
	case TypeAdmin:
		return "администратор"
	case TypeCreator:
		return "создатель"
	case TypeStudent:
		return "студент"
	case TypeOrganization:
		return "организация"
	default:
		return ""
	}
}

func (t Type) String() string {
	return string(t)
}
