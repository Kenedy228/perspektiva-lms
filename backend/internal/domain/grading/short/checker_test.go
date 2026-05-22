package short_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/short"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check_WithNormalization(t *testing.T) {
	t1, err := title.New("Вопрос")
	require.NoError(t, err)
	v := mustVariant(t, "Москва")
	q, err := qshort.New(t1, []variant.Variant{v})
	require.NoError(t, err)
	ans, err := answer.New("  москва  ")
	require.NoError(t, err)

	got, err := short.New(short.TrimSpace(), short.ToLower()).Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 1.0, got.Value())
}

func mustVariant(t *testing.T, value string) variant.Variant {
	t.Helper()

	txt, err := text.New(value)
	require.NoError(t, err)
	v, err := variant.New(txt)
	require.NoError(t, err)
	return v
}
