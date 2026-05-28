package element

type CompletionMode string

const (
	CompletionModeNone   CompletionMode = "none"
	CompletionModeManual CompletionMode = "manual"
)

func (m CompletionMode) Title() string {
	switch m {
	case CompletionModeNone:
		return "без отслеживания"
	case CompletionModeManual:
		return "ручное подтверждение"
	default:
		return ""
	}
}

func (m CompletionMode) String() string {
	return string(m)
}
