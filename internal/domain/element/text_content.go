package element

type TextContent struct {
	text string
}

func NewTextContent(text string) TextContent {
	return TextContent{
		text: text,
	}
}

func (t TextContent) Text() string {
	return t.text
}

func (t TextContent) Type() Type {
	return TypeText
}

func (t TextContent) IsInteractive() bool {
	return false
}

func (t TextContent) Clone() Content {
	return TextContent{
		text: t.text,
	}
}
