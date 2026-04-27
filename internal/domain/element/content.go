package element

type Content interface {
	Type() Type
	IsInteractive() bool
	Clone() Content
}
