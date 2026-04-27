package enrollment

type Status string

const (
	StatusInactive Status = "inactive"
	StatusActive   Status = "active"
	StatusExpired  Status = "expired"
)

func (s Status) Title() string {
	switch s {
	case StatusInactive:
		return "неактивна"
	case StatusActive:
		return "активна"
	case StatusExpired:
		return "просрочена"
	default:
		return ""
	}
}

func (s Status) String() string {
	return string(s)
}
