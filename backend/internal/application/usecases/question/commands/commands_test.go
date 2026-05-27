package commands_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/question/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase_Execute(t *testing.T) {
	t.Run("успешно создает вопрос и возвращает ID", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)

		r.On("Save", mock.Anything, mock.MatchedBy(func(q question.Question) bool {
			return q.Type() == question.TypeSelectable
		})).Return(nil).Once()

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewCreator(),
			Type:      question.TypeSelectable.String(),
			Title:     "Вопрос",
			SelectableOptions: []commands.SelectableOptionInput{
				{Text: "верно", IsCorrect: true},
				{Text: "неверно", IsCorrect: false},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.NotEmpty(t, out.ID)
		r.AssertExpectations(t)
	})

	tests := []struct {
		name      string
		in        commands.CreateInput
		setupMock func(*mockRepository)
		assertErr func(*testing.T, error)
	}{
		{
			name: "ошибка роли",
			in: commands.CreateInput{
				ActorRole: role.NewStudent(),
				Type:      question.TypeSelectable.String(),
				Title:     "Вопрос",
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrForbidden)
			},
		},
		{
			name: "ошибка типа вопроса",
			in: commands.CreateInput{
				ActorRole: role.NewAdmin(),
				Type:      "unknown",
				Title:     "Вопрос",
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrInvalidInput)
			},
		},
		{
			name: "ошибка создания заголовка",
			in: commands.CreateInput{
				ActorRole: role.NewAdmin(),
				Type:      question.TypeShort.String(),
				Title:     "",
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "ошибка создания доменного агрегата",
			in: commands.CreateInput{
				ActorRole: role.NewAdmin(),
				Type:      question.TypeSelectable.String(),
				Title:     "Вопрос",
				SelectableOptions: []commands.SelectableOptionInput{
					{Text: "нет корректного", IsCorrect: false},
					{Text: "тоже нет", IsCorrect: false},
				},
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "ошибка сохранения вопроса",
			in: commands.CreateInput{
				ActorRole: role.NewAdmin(),
				Type:      question.TypeShort.String(),
				Title:     "Вопрос",
				ShortVariants: []commands.ShortVariantInput{
					{Text: "ответ"},
				},
			},
			setupMock: func(r *mockRepository) {
				r.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			assertErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "сохранение вопроса")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mockRepository{}
			uc := commands.NewCreateUseCase(&r)
			if tt.setupMock != nil {
				tt.setupMock(&r)
			}

			out, err := uc.Execute(context.Background(), tt.in)

			assert.Nil(t, out)
			tt.assertErr(t, err)
			if tt.setupMock == nil {
				r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
			} else {
				r.AssertExpectations(t)
			}
		})
	}
}

func TestChangeTitleUseCase_Execute(t *testing.T) {
	t.Run("успешно меняет заголовок", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewChangeTitleUseCase(&r)
		q := selectableFixture(t)

		r.On("FindByID", mock.Anything, q.ID()).Return(q, nil).Once()
		r.On("Save", mock.Anything, q).Return(nil).Once()

		out, err := uc.Execute(context.Background(), commands.ChangeTitleInput{
			ActorRole:  role.NewAdmin(),
			QuestionID: q.ID().String(),
			Title:      "Новый вопрос",
		})

		require.NoError(t, err)
		require.NotNil(t, out)
		assert.Equal(t, q.ID().String(), out.ID)
		assert.Equal(t, "Новый вопрос", q.Title().Value())
		r.AssertExpectations(t)
	})

	t.Run("ошибка загрузки вопроса", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewChangeTitleUseCase(&r)
		qid := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

		r.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("not found")).Once()

		out, err := uc.Execute(context.Background(), commands.ChangeTitleInput{
			ActorRole:  role.NewCreator(),
			QuestionID: qid,
			Title:      "Новый вопрос",
		})

		assert.Nil(t, out)
		require.Error(t, err)
		assert.ErrorContains(t, err, "загрузка вопроса")
		r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		r.AssertExpectations(t)
	})

	t.Run("ошибка сохранения вопроса", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewChangeTitleUseCase(&r)
		q := selectableFixture(t)

		r.On("FindByID", mock.Anything, q.ID()).Return(q, nil).Once()
		r.On("Save", mock.Anything, q).Return(errors.New("db error")).Once()

		out, err := uc.Execute(context.Background(), commands.ChangeTitleInput{
			ActorRole:  role.NewCreator(),
			QuestionID: q.ID().String(),
			Title:      "Новый вопрос",
		})

		assert.Nil(t, out)
		require.Error(t, err)
		assert.ErrorContains(t, err, "сохранение вопроса")
		r.AssertExpectations(t)
	})
}

func TestChangeMatchingPairsUseCase_Execute_RejectsWrongType(t *testing.T) {
	r := mockRepository{}
	uc := commands.NewChangeMatchingPairsUseCase(&r)
	q := selectableFixture(t)

	r.On("FindByID", mock.Anything, q.ID()).Return(q, nil).Once()

	out, err := uc.Execute(context.Background(), commands.ChangeMatchingPairsInput{
		ActorRole:  role.NewAdmin(),
		QuestionID: q.ID().String(),
		Pairs: []commands.MatchingPairInput{
			{Prompt: "A", Match: "1"},
			{Prompt: "B", Match: "2"},
		},
	})

	assert.Nil(t, out)
	require.Error(t, err)
	assert.ErrorIs(t, err, common.ErrInvalidInput)
	r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	r.AssertExpectations(t)
}

func selectableFixture(t *testing.T) *qselectable.Question {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)
	b, err := base.New(ttl)
	require.NoError(t, err)
	opt1 := optionFixture(t, "верно", true)
	opt2 := optionFixture(t, "неверно", false)
	q, err := qselectable.New(b, []option.Option{opt1, opt2})
	require.NoError(t, err)
	return q
}

func optionFixture(t *testing.T, value string, correct bool) option.Option {
	t.Helper()

	opt, err := option.New(value, correct)
	require.NoError(t, err)
	return opt
}
