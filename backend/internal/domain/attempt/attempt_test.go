package attempt

import (
	"errors"
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"gitflic.ru/lms/backend/internal/domain/question"
	qtitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"github.com/google/uuid"
)

type mockQuestion struct {
	id    uuid.UUID
	qType question.Type
}

func (q mockQuestion) ID() uuid.UUID                  { return q.id }
func (q mockQuestion) Instruction() string            { return q.qType.DefaultInstruction() }
func (q mockQuestion) Type() question.Type            { return q.qType }
func (q mockQuestion) Clone() question.Question       { return mockQuestion{id: q.id, qType: q.qType} }
func (q mockQuestion) ChangeTitle(qtitle.Title) error { return nil }
func (q mockQuestion) Title() qtitle.Title            { t, _ := qtitle.New("Вопрос"); return t }
func asQuestions(items ...mockQuestion) []question.Question {
	qs := make([]question.Question, 0, len(items))
	for i := range items {
		qs = append(qs, items[i])
	}
	return qs
}

func mustTimeLimit(t *testing.T, seconds int) limit.Time {
	t.Helper()

	timeLimit, err := limit.NewTime(seconds)
	if err != nil {
		t.Fatalf("create time limit: %v", err)
	}

	return timeLimit
}

func validParams(t *testing.T, questions ...mockQuestion) Params {
	t.Helper()

	return Params{
		EnrollmentID: uuid.New(),
		QuizID:       uuid.New(),
		TimeLimit:    mustTimeLimit(t, 0),
		Questions:    asQuestions(questions...),
	}
}

func validAnswer(t *testing.T) question.Answer {
	t.Helper()

	return shortanswer.New("ответ")
}

func TestNew(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}

	tests := []struct {
		name    string
		params  Params
		at      time.Time
		wantErr error
	}{
		{name: "ok", params: validParams(t, q), at: now},
		{name: "zero started at", params: validParams(t, q), at: time.Time{}, wantErr: ErrInvalid},
		{name: "empty questions", params: validParams(t), at: now, wantErr: ErrInvalid},
		{name: "duplicate questions", params: validParams(t, q, q), at: now, wantErr: ErrInvalid},
		{name: "nil question", params: Params{EnrollmentID: uuid.New(), QuizID: uuid.New(), TimeLimit: mustTimeLimit(t, 0), Questions: []question.Question{nil}}, at: now, wantErr: ErrInvalid},
		{name: "zero enrollment", params: Params{EnrollmentID: uuid.Nil, QuizID: uuid.New(), TimeLimit: mustTimeLimit(t, 0), Questions: asQuestions(q)}, at: now, wantErr: ErrInvalid},
		{name: "zero quiz", params: Params{EnrollmentID: uuid.New(), QuizID: uuid.Nil, TimeLimit: mustTimeLimit(t, 0), Questions: asQuestions(q)}, at: now, wantErr: ErrInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.params, tt.at)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Status() != StatusInProgress {
				t.Fatalf("unexpected status: %s", got.Status())
			}
			if got.CountItems() != 1 {
				t.Fatalf("unexpected items count: %d", got.CountItems())
			}
		})
	}
}

func TestRestore(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}
	ans := validAnswer(t)
	entry, err := answer.New(q.ID(), ans, now.Add(time.Second))
	if err != nil {
		t.Fatalf("create entry: %v", err)
	}

	tests := []struct {
		name    string
		id      uuid.UUID
		params  RestoreParams
		wantErr error
	}{
		{
			name:   "ok in progress",
			id:     uuid.New(),
			params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusInProgress, StartedAt: now, Questions: asQuestions(q)},
		},
		{
			name:   "ok finished with answer",
			id:     uuid.New(),
			params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusFinished, StartedAt: now, FinishedAt: now.Add(time.Minute), Questions: asQuestions(q), Answers: map[uuid.UUID]answer.Entry{q.ID(): entry}},
		},
		{name: "zero id", id: uuid.Nil, params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusInProgress, StartedAt: now, Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "unknown status", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: Status("oops"), StartedAt: now, Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "finished without finishedAt", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusFinished, StartedAt: now, Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "in progress with finishedAt", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusInProgress, StartedAt: now, FinishedAt: now.Add(time.Second), Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "deadline before started", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusInProgress, StartedAt: now, DeadlineAt: now.Add(-time.Second), Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "finished before started", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusFinished, StartedAt: now, FinishedAt: now.Add(-time.Second), Questions: asQuestions(q)}, wantErr: ErrInvalid},
		{name: "answer question not found", id: uuid.New(), params: RestoreParams{EnrollmentID: uuid.New(), QuizID: uuid.New(), Status: StatusInProgress, StartedAt: now, Questions: asQuestions(q), Answers: map[uuid.UUID]answer.Entry{uuid.New(): entry}}, wantErr: ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.id, tt.params)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID() != tt.id {
				t.Fatalf("unexpected id: %s", got.ID())
			}
		})
	}
}

func TestAddAnswer(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}
	base, err := New(validParams(t, q), now)
	if err != nil {
		t.Fatalf("create attempt: %v", err)
	}

	tests := []struct {
		name    string
		mutate  func(a *Attempt)
		qID     uuid.UUID
		ans     question.Answer
		at      time.Time
		wantErr error
	}{
		{name: "ok", mutate: func(a *Attempt) {}, qID: q.ID(), ans: validAnswer(t), at: now.Add(time.Second)},
		{name: "unknown question", mutate: func(a *Attempt) {}, qID: uuid.New(), ans: validAnswer(t), at: now.Add(time.Second), wantErr: ErrNotFound},
		{name: "nil answer", mutate: func(a *Attempt) {}, qID: q.ID(), ans: nil, at: now.Add(time.Second), wantErr: ErrInvalid},
		{name: "late answer", mutate: func(a *Attempt) { a.deadlineAt = now.Add(time.Second) }, qID: q.ID(), ans: validAnswer(t), at: now.Add(2 * time.Second), wantErr: ErrStateConflict},
		{name: "cannot modify finished", mutate: func(a *Attempt) { a.status = StatusFinished; a.finishedAt = now.Add(time.Second) }, qID: q.ID(), ans: validAnswer(t), at: now.Add(time.Second), wantErr: ErrStateConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := *base
			a.answers = make(map[uuid.UUID]answer.Entry)
			tt.mutate(&a)

			err := a.AddAnswer(tt.qID, tt.ans, tt.at)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if a.CountAnswers() != 1 {
				t.Fatalf("unexpected answers count: %d", a.CountAnswers())
			}
		})
	}
}

func TestLifecycle(t *testing.T) {
	now := time.Now()
	q := mockQuestion{id: uuid.New(), qType: question.TypeShort}

	tests := []struct {
		name       string
		call       func(a *Attempt, at time.Time) error
		at         time.Time
		prepare    func(a *Attempt)
		wantStatus Status
		wantErr    error
	}{
		{name: "finish ok", call: (*Attempt).Finish, at: now.Add(time.Second), prepare: func(a *Attempt) {}, wantStatus: StatusFinished},
		{name: "finish after deadline", call: (*Attempt).Finish, at: now.Add(3 * time.Second), prepare: func(a *Attempt) { a.deadlineAt = now.Add(2 * time.Second) }, wantErr: ErrStateConflict},
		{name: "set expired ok", call: (*Attempt).SetExpired, at: now.Add(3 * time.Second), prepare: func(a *Attempt) { a.deadlineAt = now.Add(2 * time.Second) }, wantStatus: StatusExpired},
		{name: "set expired without deadline", call: (*Attempt).SetExpired, at: now.Add(time.Second), prepare: func(a *Attempt) {}, wantErr: ErrStateConflict},
		{name: "interrupt ok", call: (*Attempt).Interrupt, at: now.Add(time.Second), prepare: func(a *Attempt) {}, wantStatus: StatusInterrupted},
		{name: "cancel ok", call: (*Attempt).Cancel, at: now.Add(time.Second), prepare: func(a *Attempt) {}, wantStatus: StatusCancelled},
		{name: "cancel in final state", call: (*Attempt).Cancel, at: now.Add(time.Second), prepare: func(a *Attempt) { a.status = StatusFinished; a.finishedAt = now.Add(time.Second) }, wantErr: ErrStateConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := New(validParams(t, q), now)
			if err != nil {
				t.Fatalf("create attempt: %v", err)
			}
			tt.prepare(a)

			err = tt.call(a, tt.at)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if a.Status() != tt.wantStatus {
				t.Fatalf("unexpected status: %s", a.Status())
			}
			if !a.FinishedAt().Equal(tt.at) {
				t.Fatalf("unexpected finishedAt: %v", a.FinishedAt())
			}
		})
	}
}
