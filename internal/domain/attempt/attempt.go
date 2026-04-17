package attempt

import (
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	id           uuid.UUID
	enrollmentID uuid.UUID
	courseID     uuid.UUID
	quizID       uuid.UUID
	courseTitle  string
	quizTitle    string
	status       Status
	startedAt    time.Time
	deadlineAt   time.Time
	finishedAt   time.Time
	items        []Item
}

func NewAttempt(params Params) (*Attempt, error) {
	if err := validateEnrollmentID(params.EnrollmentID); err != nil {
		return nil, err
	}

	if err := validateCourseID(params.CourseID); err != nil {
		return nil, err
	}

	if err := validateQuizID(params.QuizID); err != nil {
		return nil, err
	}

	if err := validateCourseTitle(params.CourseTitle); err != nil {
		return nil, err
	}

	if err := validateQuizTitle(params.QuizTitle); err != nil {
		return nil, err
	}

	return nil, nil
}
