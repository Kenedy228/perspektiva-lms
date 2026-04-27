package attempt

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/shared/limit"
	"github.com/google/uuid"
)

type Params struct {
	EnrollmentID uuid.UUID
	QuizID       uuid.UUID
	Questions    []question.Question
	TimeLimit    limit.Limit
}
