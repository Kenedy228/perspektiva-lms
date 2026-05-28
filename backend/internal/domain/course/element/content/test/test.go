package test

import (
	element2 "gitflic.ru/lms/backend/internal/domain/course/element"
	"github.com/google/uuid"
)

type Content struct {
	quizID uuid.UUID
}

func New(quizID uuid.UUID) (Content, error) {
	if quizID == uuid.Nil {
		return Content{}, element2.ErrInvalid
	}
	return Content{quizID: quizID}, nil
}

func (c Content) QuizID() uuid.UUID {
	return c.quizID
}

func (c Content) ContentType() element2.ContentType {
	return element2.ContentTypeTest
}

func (c Content) IsInteractive() bool {
	return true
}

func (c Content) Clone() element2.Content {
	return c
}
