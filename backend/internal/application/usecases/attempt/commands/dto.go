package commands

import (
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type StartInput struct {
	ActorRole    role.Role
	AccountID    string
	EnrollmentID string
	QuizID       string
	Now          time.Time
}

type AddAnswerInput struct {
	ActorRole  role.Role
	AttemptID  string
	QuestionID string
	Answer     question.Answer
	AnsweredAt time.Time
}

type FinishInput struct {
	ActorRole  role.Role
	AttemptID  string
	FinishedAt time.Time
}

type ExpireInput struct {
	AttemptID string
	ExpiredAt time.Time
}

type CancelInput struct {
	ActorRole   role.Role
	AttemptID   string
	CancelledAt time.Time
}

type InterruptInput struct {
	ActorRole     role.Role
	AttemptID     string
	InterruptedAt time.Time
}

type Output struct {
	ID string
}
