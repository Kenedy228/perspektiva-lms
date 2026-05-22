package attempt

import (
	"fmt"
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

func validateEnrollmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateQuizID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateQuestions(questions []question.Question) error {
	if len(questions) == 0 {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	seen := make(map[uuid.UUID]struct{}, len(questions))
	for i, q := range questions {
		if q == nil {
			return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, i)
		}
		if q.ID() == uuid.Nil {
			return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, i)
		}
		if _, ok := seen[q.ID()]; ok {
			return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, q.ID())
		}
		seen[q.ID()] = struct{}{}
	}
	return nil
}

func validateStatus(status Status) error {
	if !status.IsValid() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateStartedAt(at time.Time) error {
	if at.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateFinishedAt(status Status, finishedAt time.Time) error {
	switch status {
	case StatusInProgress:
		if !finishedAt.IsZero() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	case StatusFinished, StatusExpired, StatusInterrupted, StatusCancelled:
		if finishedAt.IsZero() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}
	return nil
}

func validateTimeline(startedAt, deadlineAt, finishedAt time.Time) error {
	if !deadlineAt.IsZero() && deadlineAt.Before(startedAt) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	if !finishedAt.IsZero() && finishedAt.Before(startedAt) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}
