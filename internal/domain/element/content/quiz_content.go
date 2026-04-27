package content

import (
	"gitflic.ru/lms/internal/domain/element"
	"github.com/google/uuid"
)

type QuizContent struct {
	quizID uuid.UUID
}

func NewQuizContent(quizID uuid.UUID) (QuizContent, error) {
	if err := validateQuizContent(quizID); err != nil {
		return QuizContent{}, err
	}

	return QuizContent{
		quizID: quizID,
	}, nil
}

func (c QuizContent) QuizID() uuid.UUID {
	return c.quizID
}

func (c QuizContent) Type() element.Type {
	return element.TypeQuiz
}

func (c QuizContent) IsInteractive() bool {
	return true
}

func (c QuizContent) Clone() element.Content {
	return QuizContent{
		quizID: c.quizID,
	}
}
