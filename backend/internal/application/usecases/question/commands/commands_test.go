package commands_test

import (
	"context"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/question/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase(t *testing.T) {
	t.Run("forbidden for student", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewStudent(),
			Type:      question.TypeSelectable.String(),
			Title:     "Вопрос",
		})

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})

	t.Run("creates selectable question", func(t *testing.T) {
		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)
		r.On("Save", mock.Anything, mock.MatchedBy(func(q question.Question) bool {
			return q.Type() == question.TypeSelectable
		})).Return(nil)

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
		assert.NotEmpty(t, out.ID)
		r.AssertExpectations(t)
	})
}

func TestChangeTitleUseCase(t *testing.T) {
	r := mockRepository{}
	uc := commands.NewChangeTitleUseCase(&r)
	q := selectableFixture(t)

	r.On("FindByID", mock.Anything, q.ID()).Return(q, nil)
	r.On("Save", mock.Anything, q).Return(nil)

	out, err := uc.Execute(context.Background(), commands.ChangeTitleInput{
		ActorRole:  role.NewAdmin(),
		QuestionID: q.ID().String(),
		Title:      "Новый вопрос",
	})

	require.NoError(t, err)
	assert.Equal(t, q.ID().String(), out.ID)
	assert.Equal(t, "Новый вопрос", q.Title().Value())
}

func TestChangeSelectableOptionsUseCase(t *testing.T) {
	r := mockRepository{}
	uc := commands.NewChangeSelectableOptionsUseCase(&r)
	q := selectableFixture(t)

	r.On("FindByID", mock.Anything, q.ID()).Return(q, nil)
	r.On("Save", mock.Anything, q).Return(nil)

	out, err := uc.Execute(context.Background(), commands.ChangeSelectableOptionsInput{
		ActorRole:  role.NewCreator(),
		QuestionID: q.ID().String(),
		Options: []commands.SelectableOptionInput{
			{Text: "да", IsCorrect: true},
			{Text: "нет", IsCorrect: false},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, q.ID().String(), out.ID)
	assert.Equal(t, 2, len(q.Options()))
}

func selectableFixture(t *testing.T) *qselectable.Question {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)
	opt1 := optionFixture(t, "верно", true)
	opt2 := optionFixture(t, "неверно", false)
	q, err := qselectable.New(ttl, []option.Option{opt1, opt2})
	require.NoError(t, err)
	return q
}

func optionFixture(t *testing.T, value string, correct bool) option.Option {
	t.Helper()

	txt, err := text.New(value)
	require.NoError(t, err)
	opt, err := option.New(txt, correct)
	require.NoError(t, err)
	return opt
}
