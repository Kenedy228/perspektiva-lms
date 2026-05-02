package inn

type Type string

var (
	TypeIP           Type = "ip"
	TypePhysical     Type = "physical"
	TypeOrganization Type = "organization"
)

func (t Type) Title() string {
	switch t {
	case TypeIP:
		return "ИП"
	case TypePhysical:
		return "физическое лицо"
	case TypeOrganization:
		return "юридическое лицо"
	default:
		return ""
	}
}

func (t Type) CodeLength() int {
	switch t {
	case TypeIP:
		return 12
	case TypePhysical:
		return 12
	case TypeOrganization:
		return 10
	default:
		return 0
	}
}

func (t Type) Coefficients() [][]int {
	switch t {
	case TypeIP, TypePhysical:
		return [][]int{
			{7, 2, 4, 10, 3, 5, 9, 4, 6, 8},
			{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8},
		}
	case TypeOrganization:
		return [][]int{
			{2, 4, 10, 3, 5, 9, 4, 6, 8},
		}
	default:
		return [][]int{}
	}
}

func (t Type) String() string {
	return string(t)
}
