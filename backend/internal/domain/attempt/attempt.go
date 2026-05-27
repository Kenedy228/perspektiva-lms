package attempt

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"gitflic.ru/lms/backend/internal/domain/attempt/item"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
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
	answers      map[uuid.UUID]answer.Entry
	items        []item.Item
}

func New(params Params, at time.Time) (*Attempt, error) {
	if err := validateStartedAt(at); err != nil {
		return nil, err
	}

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

	var deadlineAt time.Time
	if !params.TimeLimit.IsInfinite() {
		d, ok := params.TimeLimit.TryDuration()
		if !ok {
			return nil, fmt.Errorf("%w: не удалось вычислить ограничение времени", ErrInvalid)
		}
		deadlineAt = at.Add(d)
	}

	items := make([]item.Item, 0, len(params.Questions))
	for i := range params.Questions {
		itm, err := item.New(params.Questions[i])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
		}
		items = append(items, itm)
	}

	return &Attempt{
		id:           id,
		enrollmentID: params.EnrollmentID,
		quizID:       params.QuizID,
		status:       StatusInProgress,
		startedAt:    at,
		deadlineAt:   deadlineAt,
		items:        items,
		answers:      make(map[uuid.UUID]answer.Entry, len(items)),
	}, nil
}

func Restore(id uuid.UUID, params RestoreParams) (*Attempt, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("%w: идентификатор попытки обязателен", ErrInvalid)
	}

	if err := validateEnrollmentID(params.EnrollmentID); err != nil {
		return nil, err
	}

	if err := validateQuizID(params.QuizID); err != nil {
		return nil, err
	}

	if err := validateStatus(params.Status); err != nil {
		return nil, err
	}

	if err := validateStartedAt(params.StartedAt); err != nil {
		return nil, err
	}

	if err := validateFinishedAt(params.Status, params.FinishedAt); err != nil {
		return nil, err
	}

	if err := validateTimeline(params.StartedAt, params.DeadlineAt, params.FinishedAt); err != nil {
		return nil, err
	}

	if err := validateQuestions(params.Questions); err != nil {
		return nil, err
	}

	items := make([]item.Item, 0, len(params.Questions))
	for i := range params.Questions {
		itm, err := item.New(params.Questions[i])
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
		}
		items = append(items, itm)
	}

	answers := maps.Clone(params.Answers)
	for questionID, entry := range answers {
		_, ok := findItem(items, questionID)
		if !ok {
			return nil, fmt.Errorf("%w: вопрос с идентификатором %s не найден в попытке", ErrNotFound, questionID)
		}
		if entry.QuestionID() != questionID {
			return nil, fmt.Errorf(
				"%w: ключ ответа %s не совпадает с questionID записи %s",
				ErrInvalid,
				questionID,
				entry.QuestionID(),
			)
		}
	}

	return &Attempt{
		id:           id,
		enrollmentID: params.EnrollmentID,
		quizID:       params.QuizID,
		status:       params.Status,
		startedAt:    params.StartedAt,
		deadlineAt:   params.DeadlineAt,
		finishedAt:   params.FinishedAt,
		items:        items,
		answers:      answers,
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

func (a *Attempt) Items() []item.Item {
	return slices.Clone(a.items)
}

func (a *Attempt) Answers() map[uuid.UUID]answer.Entry {
	return maps.Clone(a.answers)
}

func (a *Attempt) CountItems() int {
	return len(a.items)
}

func (a *Attempt) CountAnswers() int {
	return len(a.answers)
}

func (a *Attempt) CanModify() bool {
	return a.status == StatusInProgress
}

func (a *Attempt) AddAnswer(qID uuid.UUID, ans question.Answer, at time.Time) error {
	if !a.CanModify() {
		return fmt.Errorf("%w: попытку нельзя изменить в статусе %s", ErrStateConflict, a.status)
	}

	if err := a.ensureBeforeDeadline(at); err != nil {
		return err
	}

	if _, ok := findItem(a.items, qID); !ok {
		return fmt.Errorf("%w: вопрос с идентификатором %s не найден в попытке", ErrNotFound, qID)
	}

	entry, err := answer.New(qID, ans, at)
	if err != nil {
		return fmt.Errorf("%w: не удалось добавить ответ: %w", ErrInvalid, err)
	}

	a.answers[qID] = entry
	return nil
}

func (a *Attempt) Finish(at time.Time) error {
	if a.status != StatusInProgress {
		return fmt.Errorf("%w: попытку нельзя завершить в статусе %s", ErrStateConflict, a.status)
	}

	if at.IsZero() {
		return fmt.Errorf("%w: время завершения попытки обязательно", ErrInvalid)
	}

	if at.Before(a.startedAt) {
		return fmt.Errorf("%w: время завершения не может быть раньше времени начала", ErrInvalid)
	}

	if err := a.ensureBeforeDeadline(at); err != nil {
		return err
	}

	a.status = StatusFinished
	a.finishedAt = at
	return nil
}

func (a *Attempt) SetExpired(at time.Time) error {
	if a.status != StatusInProgress {
		return fmt.Errorf("%w: попытку нельзя пометить просроченной в статусе %s", ErrStateConflict, a.status)
	}

	if at.IsZero() {
		return fmt.Errorf("%w: время установки статуса просрочено обязательно", ErrInvalid)
	}

	if a.deadlineAt.IsZero() {
		return fmt.Errorf("%w: невозможно установить просрочку без дедлайна", ErrStateConflict)
	}

	if at.Before(a.deadlineAt) {
		return fmt.Errorf("%w: попытка не достигла дедлайна", ErrStateConflict)
	}

	a.status = StatusExpired
	a.finishedAt = at
	return nil
}

func (a *Attempt) Interrupt(at time.Time) error {
	if a.status != StatusInProgress {
		return fmt.Errorf("%w: попытку нельзя прервать в статусе %s", ErrStateConflict, a.status)
	}

	if at.IsZero() {
		return fmt.Errorf("%w: время прерывания попытки обязательно", ErrInvalid)
	}

	if at.Before(a.startedAt) {
		return fmt.Errorf("%w: время прерывания не может быть раньше времени начала", ErrInvalid)
	}

	a.status = StatusInterrupted
	a.finishedAt = at
	return nil
}

func (a *Attempt) Cancel(at time.Time) error {
	if a.status != StatusInProgress {
		return fmt.Errorf("%w: попытку нельзя отменить в статусе %s", ErrStateConflict, a.status)
	}

	if at.IsZero() {
		return fmt.Errorf("%w: время отмены попытки обязательно", ErrInvalid)
	}

	if at.Before(a.startedAt) {
		return fmt.Errorf("%w: время отмены не может быть раньше времени начала", ErrInvalid)
	}

	a.status = StatusCancelled
	a.finishedAt = at
	return nil
}

func (a *Attempt) ensureBeforeDeadline(at time.Time) error {
	if a.deadlineAt.IsZero() {
		return nil
	}

	if at.After(a.deadlineAt) {
		return fmt.Errorf("%w: момент %s позже дедлайна %s", ErrStateConflict, at, a.deadlineAt)
	}

	return nil
}
