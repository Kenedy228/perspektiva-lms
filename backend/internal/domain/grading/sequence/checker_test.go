package sequence_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/sequence"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	t1, err := title.New("Вопрос")
	require.NoError(t, err)
	opt1 := mustOption(t, "первый")
	opt2 := mustOption(t, "второй")
	q, err := qsequence.New(t1, []option.Option{opt1, opt2})
	require.NoError(t, err)

	id1, err := answer.NewOptionID(opt1.ID())
	require.NoError(t, err)
	id2, err := answer.NewOptionID(opt2.ID())
	require.NoError(t, err)
	ans, err := answer.New([]answer.OptionID{id1, id2})
	require.NoError(t, err)

	got, err := sequence.New().Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 1.0, got.Value())
}

func mustOption(t *testing.T, value string) option.Option {
	t.Helper()

	txt, err := text.New(value)
	require.NoError(t, err)
	opt, err := option.New(txt)
	require.NoError(t, err)
	return opt
}
