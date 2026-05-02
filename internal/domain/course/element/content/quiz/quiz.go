package quiz

import (
	"gitflic.ru/lms/internal/domain/course/element"
	"github.com/google/uuid"
)

type Content struct {
	quizID uuid.UUID
}

func New(quizID uuid.UUID) (Content, error) {
	if err := validateQuizID(quizID); err != nil {
		return Content{}, err
	}

	return Content{
		quizID: quizID,
	}, nil
}

func (c Content) QuizID() uuid.UUID {
	return c.quizID
}

func (c Content) ContentType() element.ContentType {
	return element.ContentTypeQuiz
}

func (c Content) IsInteractive() bool {
	return true
}

func (c Content) Clone() element.Content {
	return c
}
