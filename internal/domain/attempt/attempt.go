package attempt

import (
	"fmt"
	"time"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Attempt struct {
	id           uuid.UUID
	enrollmentID uuid.UUID
	quizID       uuid.UUID
	status       Status
	startedAt    time.Time
	deadlineAt   time.Time
	finishedAt   time.Time
	items        []*Item
}

func New(params Params) (*Attempt, error) {
	if err := validateEnrollmentID(params.EnrollmentID); err != nil {
		return nil, err
	}

	if err := validateQuizID(params.QuizID); err != nil {
		return nil, err
	}

	if err := validateQuestions(params.Questions); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var deadlineAt time.Time
	if !params.TimeLimit.IsInfinite() {
		deadlineAt = now.Add(params.TimeLimit.Duration())
	}

	items := make([]*Item, 0, len(params.Questions))
	for i := range params.Questions {
		item, err := newItem(params.Questions[i])
		if err != nil {
			return nil, fmt.Errorf("error generating item")
		}
		items = append(items, item)
	}

	return &Attempt{
		id:           id,
		enrollmentID: params.EnrollmentID,
		quizID:       params.QuizID,
		status:       StatusInProgress,
		startedAt:    now,
		deadlineAt:   deadlineAt,
		items:        items,
	}, nil
}

func (a *Attempt) ID() uuid.UUID {
	return a.id
}

func (a *Attempt) EnrollmentID() uuid.UUID {
	return a.enrollmentID
}

func (a *Attempt) QuizID() uuid.UUID {
	return a.quizID
}

func (a *Attempt) Status() Status {
	return a.status
}

func (a *Attempt) StartedAt() time.Time {
	return a.startedAt
}

func (a *Attempt) DeadlineAt() time.Time {
	return a.deadlineAt
}

func (a *Attempt) FinishedAt() time.Time {
	return a.finishedAt
}

func (a *Attempt) Items() []Item {
	items := make([]Item, 0, len(a.items))

	for i := range a.items {
		items = append(items, a.items[i].Clone())
	}

	return items
}

func (a *Attempt) CountItems() int {
	return len(a.items)
}

func (a *Attempt) CanModify() bool {
	if a.status == StatusInProgress {
		return true
	}

	return false
}

func (a *Attempt) AddAnswer(itemID uuid.UUID, answer question.Answer) error {
	if !a.CanModify() {
		return ErrNotModifiable
	}

	if itemID == uuid.Nil {
		return ErrUnexistingItem
	}

	for i := range a.items {
		if a.items[i].ID() == itemID {
			a.items[i].changeAnswer(answer)
			return nil
		}
	}

	return ErrUnexistingItem
}

func (a *Attempt) Finish() error {
	if a.status != StatusInProgress {
		return ErrInactive
	}

	a.status = StatusFinished
	a.finishedAt = time.Now()

	for i := range a.items {
		a.items[i].calculateScore()
	}

	return nil
}

func (a *Attempt) SetExpired() error {
	if a.status != StatusInProgress {
		return ErrInactive
	}

	if a.deadlineAt.IsZero() {
		return ErrInfiniteDeadline
	}

	now := time.Now()
	if now.Before(a.deadlineAt) {
		return ErrNotExpiredYet
	}

	a.status = StatusExpired
	return nil
}

func (a *Attempt) Score() (int, error) {
	if a.status != StatusFinished {
		return -1, ErrNotFinishedYet
	}

	totalScore := 0

	for i := range a.items {
		totalScore += a.items[i].Score()
	}

	return totalScore, nil
}

func (a *Attempt) Cancel() {
	if a.status == StatusCancelled {
		return
	}

	a.status = StatusCancelled
}

func (a *Attempt) Interrupt() error {
	if a.status != StatusInProgress {
		return ErrInactive
	}

	a.status = StatusInterrupted
	return nil
}
