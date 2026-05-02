package attachment

import "gitflic.ru/lms/internal/domain/question/content"

type Attachment struct {
	c content.Content
}

func New(c content.Content) (Attachment, error) {
	if err := validateContent(c); err != nil {
		return Attachment{}, err
	}

	return Attachment{
		c: c,
	}, nil
}

func (a Attachment) Content() content.Content {
	return a.c
}

func (a Attachment) Value() string {
	return a.c.Value()
}
