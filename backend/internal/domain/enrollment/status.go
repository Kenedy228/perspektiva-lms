package enrollment

type Status string

const (
	StatusInactive Status = "inactive"
	StatusActive   Status = "active"
	StatusExpired  Status = "expired"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusInactive, StatusActive, StatusExpired:
		return true
	default:
		return false
	}
}

func (s Status) Title() string {
	switch s {
	case StatusActive:
		return "активен"
	case StatusInactive:
		return "неактивен"
	case StatusExpired:
		return "истек"
	default:
		return ""
	}
}

func (s Status) String() string {
	return string(s)
}
