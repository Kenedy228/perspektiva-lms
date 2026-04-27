package content

import "gitflic.ru/lms/internal/domain/element"

type TextContent struct {
	text string
}

func NewTextContent(text string) TextContent {
	return TextContent{
		text: text,
	}
}

func (c TextContent) Text() string {
	return c.text
}

func (c TextContent) Type() element.Type {
	return element.TypeText
}

func (c TextContent) IsInteractive() bool {
	return false
}

func (c TextContent) Clone() element.Content {
	return TextContent{
		text: c.text,
	}
}
