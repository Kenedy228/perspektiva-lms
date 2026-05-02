package option

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Option struct {
	id        uuid.UUID
	c         content.Content
	isCorrect bool
}

func New(c content.Content, isCorrect bool) (Option, error) {
	if err := validateContent(c); err != nil {
		return Option{}, err
	}

	id, err := uid.New()
	if err != nil {
		return Option{}, err
	}

	return Option{
		id:        id,
		c:         c,
		isCorrect: isCorrect,
	}, nil
}

func (o Option) ID() uuid.UUID {
	return o.id
}

func (o Option) Content() content.Content {
	return o.c
}

func (o Option) IsCorrect() bool {
	return o.isCorrect
}
