package grading_test

import (
	"context"
	"errors"
	"testing"

	appgrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
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
	question  question.Question
	err       error
	findCalls int
}

func (r *fakeRepository) FindByID(_ context.Context, _ uuid.UUID) (question.Question, error) {
	r.findCalls++
	if r.err != nil {
		return nil, r.err
	}
	return r.question, nil
}

func (r *fakeRepository) Save(_ context.Context, _ question.Question) error { return nil }
func (r *fakeRepository) DeleteByID(_ context.Context, _ uuid.UUID) error   { return nil }

type fakeChecker struct {
	result     score.Score
	err        error
	checkCalls int
}

func (c *fakeChecker) Supports(_ question.Type) bool { return true }

func (c *fakeChecker) Check(_ question.Question, _ question.Answer) (score.Score, error) {
	c.checkCalls++
	if c.err != nil {
		return score.Score{}, c.err
	}
	return c.result, nil
}

type fakeValidator struct {
	err error
}

func (v fakeValidator) Validate(_ question.Question, _ question.Answer) error {
	return v.err
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
		checker        *fakeChecker
		validator      fakeValidator
		registryTypes  map[question.Type]domaingrading.Checker
		validatorTypes map[question.Type]domaingrading.AnswerValidator
		wantErr        error
		wantFindCalls  int
		wantCheckCalls int
		wantScore      float64
	}{
		{
			name:           "успешная оценка",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantFindCalls:  1,
			wantCheckCalls: 1,
			wantScore:      1,
		},
		{
			name:           "пустой question id",
			input:          appgrading.GradeInput{QuestionID: "", Answer: ans},
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
		},
		{
			name:           "некорректный uuid",
			input:          appgrading.GradeInput{QuestionID: "not-uuid", Answer: ans},
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
		},
		{
			name:           "uuid nil",
			input:          appgrading.GradeInput{QuestionID: uuid.Nil.String(), Answer: ans},
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
		},
		{
			name:           "nil answer",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: nil},
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        appgrading.ErrInvalidInput,
			wantFindCalls:  0,
		},
		{
			name:           "ошибка репозитория",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			repoErr:        repoErr,
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        repoErr,
			wantFindCalls:  1,
		},
		{
			name:           "ошибка checker",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checker:        &fakeChecker{err: checkerErr},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{err: checkerErr}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        checkerErr,
			wantFindCalls:  1,
			wantCheckCalls: 1,
		},
		{
			name:           "checker не найден в registry",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{},
			registryTypes:  map[question.Type]domaingrading.Checker{},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:        registry.ErrNotFound,
			wantFindCalls:  1,
		},
		{
			name:           "ошибка валидации ответа",
			input:          appgrading.GradeInput{QuestionID: q.ID().String(), Answer: ans},
			question:       q,
			checker:        &fakeChecker{result: sOK},
			validator:      fakeValidator{err: domaingrading.ErrInvalidAnswerType},
			registryTypes:  map[question.Type]domaingrading.Checker{question.TypeSelectable: &fakeChecker{result: sOK}},
			validatorTypes: map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{err: domaingrading.ErrInvalidAnswerType}},
			wantErr:        domaingrading.ErrInvalidAnswerType,
			wantFindCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepository{question: tt.question, err: tt.repoErr}
			if repo.question == nil {
				repo.question = q
			}

			reg, err := registry.New(tt.registryTypes)
			require.NoError(t, err)

			uc := appgrading.NewGradeUseCase(repo, reg, tt.validatorTypes)
			out, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr), "expected %v, got %v", tt.wantErr, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, out)
				assert.Equal(t, tt.wantScore, out.Score.Value())
			}

			assert.Equal(t, tt.wantFindCalls, repo.findCalls)
		})
	}
}

func TestValidateAnswerUseCase_Execute(t *testing.T) {
	q, ans := selectableFixture(t)

	tests := []struct {
		name          string
		input         appgrading.ValidateAnswerInput
		repoErr       error
		validator     fakeValidator
		validatorMap  map[question.Type]domaingrading.AnswerValidator
		wantErr       error
		wantFindCalls int
	}{
		{
			name:          "успешная валидация",
			input:         appgrading.ValidateAnswerInput{QuestionID: q.ID().String(), Answer: ans},
			validator:     fakeValidator{},
			validatorMap:  map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantFindCalls: 1,
		},
		{
			name:          "ошибка валидации входа",
			input:         appgrading.ValidateAnswerInput{QuestionID: "", Answer: ans},
			validator:     fakeValidator{},
			validatorMap:  map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{}},
			wantErr:       appgrading.ErrInvalidInput,
			wantFindCalls: 0,
		},
		{
			name:          "ошибка валидации ответа",
			input:         appgrading.ValidateAnswerInput{QuestionID: q.ID().String(), Answer: ans},
			validator:     fakeValidator{err: domaingrading.ErrInvalidAnswerType},
			validatorMap:  map[question.Type]domaingrading.AnswerValidator{question.TypeSelectable: fakeValidator{err: domaingrading.ErrInvalidAnswerType}},
			wantErr:       domaingrading.ErrInvalidAnswerType,
			wantFindCalls: 1,
		},
		{
			name:          "валидатор не найден",
			input:         appgrading.ValidateAnswerInput{QuestionID: q.ID().String(), Answer: ans},
			validator:     fakeValidator{},
			validatorMap:  map[question.Type]domaingrading.AnswerValidator{question.TypeMatching: fakeValidator{}},
			wantErr:       errors.New("валидатор для типа вопроса"),
			wantFindCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepository{question: q, err: tt.repoErr}
			uc := appgrading.NewValidateAnswerUseCase(repo, tt.validatorMap)

			err := uc.Execute(context.Background(), tt.input)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantFindCalls, repo.findCalls)
		})
	}
}

func TestNewGradeUseCase_Panics(t *testing.T) {
	q, _ := selectableFixture(t)
	repo := &fakeRepository{question: q}

	reg, err := registry.New(map[question.Type]domaingrading.Checker{
		question.TypeSelectable: &fakeChecker{},
	})
	require.NoError(t, err)

	validators := map[question.Type]domaingrading.AnswerValidator{
		question.TypeSelectable: fakeValidator{},
	}

	assert.Panics(t, func() {
		appgrading.NewGradeUseCase(nil, reg, validators)
	})

	assert.Panics(t, func() {
		appgrading.NewGradeUseCase(repo, nil, validators)
	})

	assert.Panics(t, func() {
		appgrading.NewGradeUseCase(repo, reg, nil)
	})
}

func TestNewValidateAnswerUseCase_Panics(t *testing.T) {
	repo := &fakeRepository{}
	validators := map[question.Type]domaingrading.AnswerValidator{
		question.TypeSelectable: fakeValidator{},
	}

	assert.Panics(t, func() {
		appgrading.NewValidateAnswerUseCase(nil, validators)
	})

	assert.Panics(t, func() {
		appgrading.NewValidateAnswerUseCase(repo, nil)
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
