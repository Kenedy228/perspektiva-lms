package version

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
	StatusDeleted   Status = "deleted"
)

func (s Status) Title() string {
	switch s {
	case StatusDraft:
		return "черновик"
	case StatusPublished:
		return "опубликован"
	case StatusDeleted:
		return "удален"
	default:
		return ""
	}
}

func (s Status) IsValid() bool {
	switch s {
	case StatusDraft, StatusPublished, StatusDeleted:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}
