package attempt

import (
	"errors"
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"gitflic.ru/lms/backend/internal/domain/question"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"github.com/google/uuid"
)

func TestNewRejectsZeroStartAndDuplicateQuestions(t *testing.T) {
	params := validParams(t, mockQuestion{id: uuid.New(), qType: question.TypeShort})

	_, err := New(params, time.Time{})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid attempt for zero start, got %v", err)
	}

	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}
	params = validParams(t, q, q)
	_, err = New(params, time.Now())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid attempt for duplicate questions, got %v", err)
	}
}

func TestAddAnswerRejectsLateAndWrongTypeAnswers(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}
	params := validParams(t, q)
	params.TimeLimit = mustTimeLimit(t, 60)

	a, err := New(params, now)
	if err != nil {
		t.Fatalf("create attempt: %v", err)
	}

	ans, err := shortanswer.New("answer")
	if err != nil {
		t.Fatalf("create answer: %v", err)
	}

	err = a.AddAnswer(q.ID(), ans, now.Add(61*time.Second))
	if !errors.Is(err, ErrStateConflict) {
		t.Fatalf("expected state conflict for late answer, got %v", err)
	}

	wrong := mockQuestion{id: uuid.New(), qType: question.TypeSelectable}
	a, err = New(validParams(t, wrong), now)
	if err != nil {
		t.Fatalf("create attempt: %v", err)
	}
	err = a.AddAnswer(wrong.ID(), ans, now.Add(time.Second))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid attempt for wrong answer type, got %v", err)
	}
}

func TestFinishRejectsAfterDeadline(t *testing.T) {
	now := time.Now()
	params := validParams(t, mockQuestion{id: uuid.New(), qType: question.TypeShort})
	params.TimeLimit = mustTimeLimit(t, 60)

	a, err := New(params, now)
	if err != nil {
		t.Fatalf("create attempt: %v", err)
	}

	err = a.Finish(now.Add(61 * time.Second))
	if !errors.Is(err, ErrStateConflict) {
		t.Fatalf("expected state conflict for late finish, got %v", err)
	}
}

func TestCancelSetsTerminalTimestampAndCannotRepeat(t *testing.T) {
	now := time.Now()
	a, err := New(validParams(t, mockQuestion{id: uuid.New(), qType: question.TypeShort}), now)
	if err != nil {
		t.Fatalf("create attempt: %v", err)
	}

	cancelledAt := now.Add(time.Second)
	if err := a.Cancel(cancelledAt); err != nil {
		t.Fatalf("cancel attempt: %v", err)
	}
	if a.Status() != StatusCancelled {
		t.Fatalf("expected cancelled status, got %s", a.Status())
	}
	if !a.FinishedAt().Equal(cancelledAt) {
		t.Fatalf("expected cancelled timestamp %v, got %v", cancelledAt, a.FinishedAt())
	}

	err = a.Cancel(cancelledAt)
	if !errors.Is(err, ErrStateConflict) {
		t.Fatalf("expected conflict on repeated cancel, got %v", err)
	}
}

func TestRestoreValidatesStatusAndAnswers(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}
	ans, err := shortanswer.New("answer")
	if err != nil {
		t.Fatalf("create answer: %v", err)
	}
	entry, err := answerEntry(q.ID(), ans, now.Add(time.Second))
	if err != nil {
		t.Fatalf("create answer entry: %v", err)
	}

	_, err = Restore(uuid.New(), RestoreParams{
		EnrollmentID: uuid.New(),
		QuizID:       uuid.New(),
		Status:       Status("unknown"),
		StartedAt:    now,
		Questions:    asQuestions(q),
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid status error, got %v", err)
	}

	a, err := Restore(uuid.New(), RestoreParams{
		EnrollmentID: uuid.New(),
		QuizID:       uuid.New(),
		Status:       StatusInProgress,
		StartedAt:    now,
		Questions:    asQuestions(q),
		Answers:      map[uuid.UUID]answer.Entry{q.ID(): entry},
	})
	if err != nil {
		t.Fatalf("restore attempt: %v", err)
	}
	if a.CountAnswers() != 1 {
		t.Fatalf("expected restored answer, got %d", a.CountAnswers())
	}
}
