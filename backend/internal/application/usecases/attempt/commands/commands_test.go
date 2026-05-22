package commands

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/quiz"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	quiztitle "gitflic.ru/lms/backend/internal/domain/quiz/title"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	attemptinfra "gitflic.ru/lms/backend/internal/infrastructure/attempt"
	"github.com/google/uuid"
)

func TestStartAppliesQuizShufflePolicy(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	accountID := uuid.New()
	enrollmentID := uuid.New()
	bankID := uuid.New()
	first := testQuestion{id: uuid.New(), qType: question.TypeShort}
	second := testQuestion{id: uuid.New(), qType: question.TypeShort}
	third := testQuestion{id: uuid.New(), qType: question.TypeShort}
	fourth := testQuestion{id: uuid.New(), qType: question.TypeShort}

	q := mustQuiz(t, bankID, true, 0, 4)
	attempts := newAttemptRepo()
	quizzes := newQuizRepo(q)
	enrollment := &enrollmentPolicy{allowed: true}
	provider := &questionProvider{
		random: map[uuid.UUID][]question.Question{bankID: {first, second, third, fourth}},
	}
	shuffler := attemptinfra.NewMathRandQuestionShufflerWithSource(rand.NewSource(1))

	out, err := NewStartUseCase(attempts, quizzes, enrollment, provider, shuffler).Execute(ctx, StartInput{
		ActorRole:    role.NewStudent(),
		AccountID:    accountID.String(),
		EnrollmentID: enrollmentID.String(),
		QuizID:       q.ID().String(),
		Now:          now,
	})
	if err != nil {
		t.Fatalf("start attempt: %v", err)
	}
	if out.ID == "" {
		t.Fatal("expected output id")
	}

	items := attempts.saved.Items()
	if items[0].ID() == first.ID() && items[1].ID() == second.ID() && items[2].ID() == third.ID() && items[3].ID() == fourth.ID() {
		t.Fatal("expected shuffled question order")
	}
}

func TestStartRejectsNonStudentAndNotEnrolled(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	bankID := uuid.New()
	q := mustQuiz(t, bankID, false, 0, 1)
	attempts := newAttemptRepo()
	quizzes := newQuizRepo(q)
	provider := &questionProvider{
		random: map[uuid.UUID][]question.Question{bankID: {testQuestion{id: uuid.New(), qType: question.TypeShort}}},
	}

	_, err := NewStartUseCase(attempts, quizzes, &enrollmentPolicy{allowed: true}, provider, reverseShuffler{}).Execute(ctx, StartInput{
		ActorRole:    role.NewCreator(),
		AccountID:    uuid.New().String(),
		EnrollmentID: uuid.New().String(),
		QuizID:       q.ID().String(),
		Now:          now,
	})
	if !errors.Is(err, common.ErrForbidden) {
		t.Fatalf("expected forbidden for non-student, got %v", err)
	}

	_, err = NewStartUseCase(attempts, quizzes, &enrollmentPolicy{allowed: false}, provider, reverseShuffler{}).Execute(ctx, StartInput{
		ActorRole:    role.NewStudent(),
		AccountID:    uuid.New().String(),
		EnrollmentID: uuid.New().String(),
		QuizID:       q.ID().String(),
		Now:          now,
	})
	if !errors.Is(err, common.ErrForbidden) {
		t.Fatalf("expected forbidden for missing enrollment, got %v", err)
	}
}

func TestStartRejectsAttemptsLimitReached(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	bankID := uuid.New()
	q := mustQuiz(t, bankID, false, 1, 1)
	attempts := newAttemptRepo()
	attempts.count = 1
	quizzes := newQuizRepo(q)
	provider := &questionProvider{
		random: map[uuid.UUID][]question.Question{bankID: {testQuestion{id: uuid.New(), qType: question.TypeShort}}},
	}

	_, err := NewStartUseCase(attempts, quizzes, &enrollmentPolicy{allowed: true}, provider, reverseShuffler{}).Execute(ctx, StartInput{
		ActorRole:    role.NewStudent(),
		AccountID:    uuid.New().String(),
		EnrollmentID: uuid.New().String(),
		QuizID:       q.ID().String(),
		Now:          now,
	})
	if !errors.Is(err, common.ErrLimitReached) {
		t.Fatalf("expected limit reached, got %v", err)
	}
}

type attemptRepo struct {
	items map[uuid.UUID]*attemptdomain.Attempt
	saved *attemptdomain.Attempt
	count int
}

func newAttemptRepo() *attemptRepo {
	return &attemptRepo{items: make(map[uuid.UUID]*attemptdomain.Attempt)}
}

func (r *attemptRepo) FindByID(_ context.Context, id uuid.UUID) (*attemptdomain.Attempt, error) {
	a, ok := r.items[id]
	if !ok {
		return nil, errors.New("attempt not found")
	}
	return a, nil
}

func (r *attemptRepo) Save(_ context.Context, a *attemptdomain.Attempt) error {
	r.items[a.ID()] = a
	r.saved = a
	return nil
}

func (r *attemptRepo) CountByEnrollmentAndQuiz(_ context.Context, _, _ uuid.UUID) (int, error) {
	return r.count, nil
}

type quizRepo struct {
	q *quiz.Quiz
}

func newQuizRepo(q *quiz.Quiz) *quizRepo {
	return &quizRepo{q: q}
}

func (r *quizRepo) FindByID(_ context.Context, id uuid.UUID) (*quiz.Quiz, error) {
	if r.q == nil || r.q.ID() != id {
		return nil, errors.New("quiz not found")
	}
	return r.q, nil
}

func (r *quizRepo) Save(context.Context, *quiz.Quiz) error { return nil }
func (r *quizRepo) DeleteByID(context.Context, uuid.UUID) error {
	return nil
}

type enrollmentPolicy struct {
	allowed bool
}

func (p *enrollmentPolicy) CanStartQuiz(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, time.Time) (bool, error) {
	return p.allowed, nil
}

type questionProvider struct {
	manual map[uuid.UUID][]question.Question
	random map[uuid.UUID][]question.Question
}

func (p *questionProvider) FindQuestionsByIDs(_ context.Context, bankID uuid.UUID, _ []uuid.UUID) ([]question.Question, error) {
	return p.manual[bankID], nil
}

func (p *questionProvider) SelectRandomQuestions(_ context.Context, bankID uuid.UUID, _ int) ([]question.Question, error) {
	return p.random[bankID], nil
}

type reverseShuffler struct{}

func (reverseShuffler) ShuffleQuestions(questions []question.Question) []question.Question {
	shuffled := make([]question.Question, 0, len(questions))
	for i := len(questions) - 1; i >= 0; i-- {
		shuffled = append(shuffled, questions[i])
	}
	return shuffled
}

type testQuestion struct {
	id    uuid.UUID
	qType question.Type
}

func (q testQuestion) ID() uuid.UUID { return q.id }

func (q testQuestion) Title() title.Title {
	t, _ := title.New("Question")
	return t
}

func (q testQuestion) Attachment() (attachment.Attachment, bool) {
	return attachment.Attachment{}, false
}

func (q testQuestion) Instruction() string { return q.qType.DefaultInstruction() }
func (q testQuestion) Type() question.Type { return q.qType }
func (q testQuestion) Clone() question.Question {
	return testQuestion{id: q.id, qType: q.qType}
}
func (q testQuestion) ChangeTitle(title.Title)                {}
func (q testQuestion) ChangeAttachment(attachment.Attachment) {}
func (q testQuestion) RemoveAttachment()                      {}
func (q testQuestion) HasAttachment() bool                    { return false }

func mustQuiz(t *testing.T, bankID uuid.UUID, shuffle bool, maxAttempts int, questionCount int) *quiz.Quiz {
	t.Helper()

	tl, err := quiztitle.New("Quiz")
	if err != nil {
		t.Fatalf("create quiz title: %v", err)
	}
	attempts, err := limit.NewAttempts(maxAttempts)
	if err != nil {
		t.Fatalf("create attempts limit: %v", err)
	}
	c, err := criteria.NewRandom(questionCount)
	if err != nil {
		t.Fatalf("create criteria: %v", err)
	}
	s, err := source.NewSource(bankID, c)
	if err != nil {
		t.Fatalf("create source: %v", err)
	}
	q, err := quiz.New(quiz.Params{
		Title:            tl,
		MaxAttempts:      attempts,
		ShuffleQuestions: shuffle,
		Sources:          []source.Source{s},
	})
	if err != nil {
		t.Fatalf("create quiz: %v", err)
	}
	return q
}
