package element

type CompletionMode string

const (
	CompletionModeNone   CompletionMode = "none"
	CompletionModeManual CompletionMode = "manual"
)

func (m CompletionMode) String() string {
	return string(m)
}
