package attempt

import (
	"time"

	"gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"github.com/google/uuid"
)

type Params struct {
	EnrollmentID uuid.UUID
	QuizID       uuid.UUID
	TimeLimit    limit.Time
	Questions    []question.Question
}

type RestoreParams struct {
	EnrollmentID uuid.UUID
	QuizID       uuid.UUID
	Status       Status
	StartedAt    time.Time
	DeadlineAt   time.Time
	FinishedAt   time.Time
	Questions    []question.Question
	Answers      map[uuid.UUID]answer.Entry
}
