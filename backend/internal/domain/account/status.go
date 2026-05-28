package account

// Status описывает, может ли аккаунт использоваться для аутентификации.
type Status string

const (
	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
	StatusDeleted Status = "deleted"
)

// IsValid показывает, поддерживается ли в системе статус.
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusBlocked, StatusDeleted:
		return true
	default:
		return false
	}
}

// Title возвращает локализованное название конкретного статуса.
func (s Status) Title() string {
	switch s {
	case StatusActive:
		return "активный"
	case StatusBlocked:
		return "заблокирован"
	case StatusDeleted:
		return "удалён"
	default:
		return ""
	}
}

// String возвращает строковое представление статуса в системе.
func (s Status) String() string {
	return string(s)
}
