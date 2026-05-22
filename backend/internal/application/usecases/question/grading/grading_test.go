package grading_test

import (
	"context"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading/selectable"
	"gitflic.ru/lms/backend/internal/domain/question"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (question.Question, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(question.Question), args.Error(1)
}

func (m *mockRepository) Save(ctx context.Context, q question.Question) error {
	args := m.Called(ctx, q)
	return args.Error(0)
}

func (m *mockRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGradeUseCase(t *testing.T) {
	r := mockRepository{}
	q, ans := selectableFixture(t)
	uc := grading.NewGradeUseCase(&r, domaingrading.New())

	r.On("FindByID", mock.Anything, q.ID()).Return(q, nil)

	out, err := uc.Execute(context.Background(), grading.GradeInput{
		QuestionID: q.ID().String(),
		Answer:     ans,
	})

	require.NoError(t, err)
	assert.Equal(t, 1.0, out.Score.Value())
}

func TestValidateAnswerUseCase(t *testing.T) {
	r := mockRepository{}
	q, ans := selectableFixture(t)
	uc := grading.NewValidateAnswerUseCase(&r, domaingrading.New())

	r.On("FindByID", mock.Anything, q.ID()).Return(q, nil)

	err := uc.Execute(context.Background(), grading.ValidateAnswerInput{
		QuestionID: q.ID().String(),
		Answer:     ans,
	})

	require.NoError(t, err)
}

func selectableFixture(t *testing.T) (*qselectable.Question, answer.Answer) {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)
	opt1 := optionFixture(t, "верно", true)
	opt2 := optionFixture(t, "неверно", false)
	q, err := qselectable.New(ttl, []option.Option{opt1, opt2})
	require.NoError(t, err)
	id, err := answer.NewOptionID(opt1.ID())
	require.NoError(t, err)
	ans, err := answer.New([]answer.OptionID{id})
	require.NoError(t, err)
	return q, ans
}

func optionFixture(t *testing.T, value string, correct bool) option.Option {
	t.Helper()

	txt, err := text.New(value)
	require.NoError(t, err)
	opt, err := option.New(txt, correct)
	require.NoError(t, err)
	return opt
}
