package variant

import "gitflic.ru/lms/internal/domain/question/content"

type Variant struct {
	c content.Content
}

func New(c content.Content) (Variant, error) {
	if err := validateContent(c); err != nil {
		return Variant{}, err
	}

	return Variant{
		c: c,
	}, nil
}

func (v Variant) Content() content.Content {
	return v.c
}
