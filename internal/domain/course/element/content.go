package element

type Content interface {
	ContentType() ContentType
	IsInteractive() bool
	Clone() Content
}
