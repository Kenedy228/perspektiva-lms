package answer

import (
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

type Entry struct {
	questionID uuid.UUID
	answer     question.Answer
	answeredAt time.Time
}

func New(questionID uuid.UUID, ans question.Answer, answeredAt time.Time) (Entry, error) {
	if err := validateQuestionID(questionID); err != nil {
		return Entry{}, err
	}

	if err := validateAnswer(ans); err != nil {
		return Entry{}, err
	}

	if err := validateAnsweredAt(answeredAt); err != nil {
		return Entry{}, err
	}

	return Entry{
		questionID: questionID,
		answer:     ans.Clone(),
		answeredAt: answeredAt,
	}, nil
}

func (e Entry) QuestionID() uuid.UUID {
	return e.questionID
}

func (e Entry) Answer() question.Answer {
	return e.answer.Clone()
}

func (e Entry) AnsweredAt() time.Time {
	return e.answeredAt
}
