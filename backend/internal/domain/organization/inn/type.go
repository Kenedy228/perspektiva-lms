package inn

type Type string

const (
	TypeIndividualEntrepreneur Type = "individual entrepreneur"
	TypeNaturalPerson          Type = "natural person"
	TypeLegalEntity            Type = "legal entity"
)

func (t Type) IsValid() bool {
	switch t {
	case TypeIndividualEntrepreneur, TypeNaturalPerson, TypeLegalEntity:
		return true
	default:
		return false
	}
}

// Title returns the display title for the INN type.
func (t Type) Title() string {
	switch t {
	case TypeIndividualEntrepreneur:
		return "индивидуальный предприниматель"
	case TypeNaturalPerson:
		return "физическое лицо"
	case TypeLegalEntity:
		return "юридическое лицо"
	default:
		return ""
	}
}

// String returns the raw type value.
func (t Type) String() string {
	return string(t)
}
