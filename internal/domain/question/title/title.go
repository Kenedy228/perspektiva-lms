package title

import "gitflic.ru/lms/internal/domain/question/content"

type Title struct {
	c content.Content
}

func New(c content.Content) (Title, error) {
	if err := validateContent(c); err != nil {
		return Title{}, err
	}

	return Title{
		c: c,
	}, nil
}

func (t Title) Value() string {
	return t.c.Value()
}
