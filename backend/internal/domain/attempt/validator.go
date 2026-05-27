package attempt

import (
	"fmt"
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

func validateEnrollmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор записи на курс обязателен", ErrInvalid)
	}
	return nil
}

func validateQuizID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор теста обязателен", ErrInvalid)
	}
	return nil
}

func validateQuestions(questions []question.Question) error {
	if len(questions) == 0 {
		return fmt.Errorf("%w: список вопросов не должен быть пустым", ErrInvalid)
	}

	seen := make(map[uuid.UUID]struct{}, len(questions))
	for i, q := range questions {
		if q == nil {
			return fmt.Errorf("%w: вопрос с индексом %d не должен быть nil", ErrInvalid, i)
		}
		if q.ID() == uuid.Nil {
			return fmt.Errorf("%w: идентификатор вопроса с индексом %d обязателен", ErrInvalid, i)
		}
		if _, ok := seen[q.ID()]; ok {
			return fmt.Errorf("%w: вопрос с идентификатором %s уже добавлен в попытку", ErrInvalid, q.ID())
		}
		seen[q.ID()] = struct{}{}
	}
	return nil
}

func validateStatus(status Status) error {
	if !status.IsValid() {
		return fmt.Errorf("%w: неизвестный статус попытки %q", ErrInvalid, status)
	}
	return nil
}

func validateStartedAt(at time.Time) error {
	if at.IsZero() {
		return fmt.Errorf("%w: время начала попытки обязательно", ErrInvalid)
	}
	return nil
}

func validateFinishedAt(status Status, finishedAt time.Time) error {
	switch status {
	case StatusInProgress:
		if !finishedAt.IsZero() {
			return fmt.Errorf("%w: время завершения должно быть пустым для статуса %s", ErrInvalid, status)
		}
	case StatusFinished, StatusExpired, StatusInterrupted, StatusCancelled:
		if finishedAt.IsZero() {
			return fmt.Errorf("%w: время завершения обязательно для статуса %s", ErrInvalid, status)
		}
	}
	return nil
}

func validateTimeline(startedAt, deadlineAt, finishedAt time.Time) error {
	if !deadlineAt.IsZero() && deadlineAt.Before(startedAt) {
		return fmt.Errorf("%w: дедлайн не может быть раньше времени начала", ErrInvalid)
	}
	if !finishedAt.IsZero() && finishedAt.Before(startedAt) {
		return fmt.Errorf("%w: время завершения не может быть раньше времени начала", ErrInvalid)
	}
	return nil
}
