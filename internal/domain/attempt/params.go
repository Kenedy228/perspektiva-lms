package attempt

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/quiz/limit"
	"github.com/google/uuid"
)

type Params struct {
	EnrollmentID uuid.UUID
	QuizID       uuid.UUID
	TimeLimit    limit.Time
	Questions    []question.Question
}
