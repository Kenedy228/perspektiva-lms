package typed_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/typed"
	qtyped "gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check_AllOrNothing(t *testing.T) {
	t1, err := title.New("Столица: {{city}}")
	require.NoError(t, err)
	city := mustBlank(t, "{{city}}", "Москва")
	q, err := qtyped.New(t1, []blank.Blank{city})
	require.NoError(t, err)
	ans, err := answer.New([]answer.AnswerBlank{{Placeholder: "{{city}}", Variant: " москва "}})
	require.NoError(t, err)

	got, err := typed.New(typed.TrimSpace(), typed.ToLower()).Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 1.0, got.Value())
}

func TestChecker_Check_ExtraBlankScoresZero(t *testing.T) {
	t1, err := title.New("Столица: {{city}}")
	require.NoError(t, err)
	city := mustBlank(t, "{{city}}", "Москва")
	q, err := qtyped.New(t1, []blank.Blank{city})
	require.NoError(t, err)
	ans, err := answer.New([]answer.AnswerBlank{
		{Placeholder: "{{city}}", Variant: "Москва"},
		{Placeholder: "{{extra}}", Variant: "x"},
	})
	require.NoError(t, err)

	got, err := typed.New().Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 0.0, got.Value())
}

func mustBlank(t *testing.T, placeholder, variant string) blank.Blank {
	t.Helper()

	txt, err := text.New(variant)
	require.NoError(t, err)
	b, err := blank.New(placeholder, []text.Text{txt})
	require.NoError(t, err)
	return b
}
