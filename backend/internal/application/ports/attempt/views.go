package attempt

import "time"

type AttemptSummaryView struct {
	ID             string
	EnrollmentID   string
	QuizID         string
	Status         string
	StartedAt      time.Time
	DeadlineAt     time.Time
	FinishedAt     time.Time
	QuestionsCount int
	AnswersCount   int
}
