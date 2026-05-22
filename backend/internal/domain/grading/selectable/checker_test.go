package selectable_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/selectable"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	t1, err := title.New("Вопрос")
	require.NoError(t, err)
	opt1 := mustOption(t, "верно", true)
	opt2 := mustOption(t, "неверно", false)
	q, err := qselectable.New(t1, []option.Option{opt1, opt2})
	require.NoError(t, err)

	id, err := answer.NewOptionID(opt1.ID())
	require.NoError(t, err)
	ans, err := answer.New([]answer.OptionID{id})
	require.NoError(t, err)

	got, err := selectable.New().Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 1.0, got.Value())
}

func mustOption(t *testing.T, value string, correct bool) option.Option {
	t.Helper()

	txt, err := text.New(value)
	require.NoError(t, err)
	opt, err := option.New(txt, correct)
	require.NoError(t, err)
	return opt
}
