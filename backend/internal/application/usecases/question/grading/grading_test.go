package grading_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	appgrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	question   question.Question
	err        error
	findCalls  int
	lastFindID uuid.UUID
}

func (r *fakeRepository) FindByID(_ context.Context, id uuid.UUID) (question.Question, error) {
	r.findCalls++
	r.lastFindID = id
	if r.err != nil {
		return nil, r.err
	}
	return r.question, nil
}

func (r *fakeRepository) Save(_ context.Context, _ question.Question) error { return nil }
func (r *fakeRepository) DeleteByID(_ context.Context, _ uuid.UUID) error   { return nil }

type fakeChecker struct {
	supportedTypes map[question.Type]bool
	result         score.Score
	err            error
	checkCalls     int
}

func (c *fakeChecker) Supports(t question.Type) bool {
	return c.supportedTypes[t]
}

func (c *fakeChecker) Check(_ question.Question, _ question.Answer) (score.Score, error) {
	c.checkCalls++
	if c.err != nil {
		return score.Score{}, c.err
	}
	return c.result, nil
}

func TestGradeUseCase_Execute(t *testing.T) {
	q, ans := selectableFixture(t)
	sOK, err := score.New(1)
	require.NoError(t, err)

	repoErr := errors.New("repo failed")
	checkerErr := errors.New("checker failed")

	tests := []struct {
		name           string
		input          appgrading.GradeInput
		repoErr        error
		question       question.Question
		checkers       []*fakeChecker
		wantErr        error
		wantFindCalls  int
		wantCheckCalls []int
		wantScore      float64
	}{
		{
			name:           "успешная оценка",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantFindCalls:  1,
			wantCheckCalls: []int{1},
			wantScore:      1,
		},
		{
			name:           "пустой question id",
			input:          appgrading.GradeInput{QuestionID: "", Answer: ans},
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
			wantCheckCalls: []int{0},
		},
		{
			name:           "некорректный uuid",
			input:          appgrading.GradeInput{QuestionID: "not-uuid", Answer: ans},
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
			wantCheckCalls: []int{0},
		},
		{
			name:           "uuid nil",
			input:          appgrading.GradeInput{QuestionID: uuid.Nil.String(), Answer: ans},
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
			wantCheckCalls: []int{0},
		},
		{
			name:           "nil answer",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: nil},
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
			wantCheckCalls: []int{0},
		},
		{
			name:           "ошибка репозитория",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			repoErr:        repoErr,
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK}},
			wantErr:        repoErr,
			wantFindCalls:  1,
			wantCheckCalls: []int{0},
		},
		{
			name:           "ошибка checker",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, err: checkerErr}},
			wantErr:        checkerErr,
			wantFindCalls:  1,
			wantCheckCalls: []int{1},
		},
		{
			name:           "поддерживаемый checker не найден",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checkers:       []*fakeChecker{{supportedTypes: map[question.Type]bool{question.TypeShort: true}, result: sOK}},
			wantErr:        appgrading.ErrUnsupportedChecker,
			wantFindCalls:  1,
			wantCheckCalls: []int{0},
		},
		{
			name:     "выбирается корректный checker по типу",
			input:    appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question: q,
			checkers: []*fakeChecker{
				{supportedTypes: map[question.Type]bool{question.TypeShort: true}, result: sOK},
				{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK},
			},
			wantFindCalls:  1,
			wantCheckCalls: []int{0, 1},
			wantScore:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepository{question: tt.question, err: tt.repoErr}
			if repo.question == nil {
				repo.question = q
			}

			checkers := make([]domaingrading.Checker, 0, len(tt.checkers))
			for i := range tt.checkers {
				checkers = append(checkers, tt.checkers[i])
			}

			uc := appgrading.NewGradeUseCase(repo, checkers...)
			out, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				require.NoError(t, err)
				require.NotNil(t, out)
				assert.Equal(t, tt.wantScore, out.Score.Value())
			}

			assert.Equal(t, tt.wantFindCalls, repo.findCalls)
			for i := range tt.checkers {
				assert.Equal(t, tt.wantCheckCalls[i], tt.checkers[i].checkCalls)
			}
		})
	}
}

func TestValidateAnswerUseCase_Execute(t *testing.T) {
	q, ans := selectableFixture(t)
	sOK, err := score.New(1)
	require.NoError(t, err)

	tests := []struct {
		name          string
		input         appgrading.ValidateAnswerInput
		repoErr       error
		checkerErr    error
		wantErr       error
		wantFindCalls int
	}{
		{
			name:          "успешная валидация без двойной загрузки",
			input:         appgrading.ValidateAnswerInput{QuestionID: q.ID().String(), Answer: ans},
			wantFindCalls: 1,
		},
		{
			name:          "ошибка валидации входа",
			input:         appgrading.ValidateAnswerInput{QuestionID: "", Answer: ans},
			wantErr:       appgrading.ErrInvalidInput,
			wantFindCalls: 0,
		},
		{
			name:          "ошибка checker",
			input:         appgrading.ValidateAnswerInput{QuestionID: q.ID().String(), Answer: ans},
			checkerErr:    errors.New("checker failed"),
			wantErr:       errors.New("checker failed"),
			wantFindCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepository{question: q, err: tt.repoErr}
			checker := &fakeChecker{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}, result: sOK, err: tt.checkerErr}
			uc := appgrading.NewValidateAnswerUseCase(repo, checker)

			err := uc.Execute(context.Background(), tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				if tt.name == "ошибка checker" {
					assert.Contains(t, err.Error(), "checker failed")
				} else {
					assert.True(t, errors.Is(err, tt.wantErr))
				}
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantFindCalls, repo.findCalls)
		})
	}
}

func TestNewGradeUseCase_Panics(t *testing.T) {
	q, _ := selectableFixture(t)
	checker := &fakeChecker{supportedTypes: map[question.Type]bool{question.TypeSelectable: true}}

	assert.Panics(t, func() {
		appgrading.NewGradeUseCase(nil, checker)
	})

	repo := &fakeRepository{question: q}
	assert.Panics(t, func() {
		appgrading.NewGradeUseCase(repo)
	})
}

func selectableFixture(t *testing.T) (*qselectable.Question, selectableanswer.Answer) {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)
	b, err := base.New(ttl)
	require.NoError(t, err)

	opt1, err := option.New("верно", true)
	require.NoError(t, err)
	opt2, err := option.New("неверно", false)
	require.NoError(t, err)

	q, err := qselectable.New(b, []option.Option{opt1, opt2})
	require.NoError(t, err)

	a, err := selectableanswer.New([]uuid.UUID{opt1.ID()})
	require.NoError(t, err)

	return q, a
}

func TestSentinelIdentity(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", appgrading.ErrUnsupportedChecker)
	assert.True(t, errors.Is(err, appgrading.ErrUnsupportedChecker))
}
