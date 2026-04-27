package version

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPublished Status = "published"
)

func (s Status) Title() string {
	switch s {
	case StatusDraft:
		return "черновик"
	case StatusPublished:
		return "опубликован"
	default:
		return ""
	}
}

func (s Status) String() string {
	return string(s)
}
