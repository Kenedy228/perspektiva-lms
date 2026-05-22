package answer_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ans, err := answer.New([]answer.AnswerBlank{{Placeholder: "{{city}}", Variant: "Москва"}})
	require.NoError(t, err)
	assert.Equal(t, []answer.AnswerBlank{{Placeholder: "{{city}}", Variant: "Москва"}}, ans.Blanks())

	_, err = answer.New([]answer.AnswerBlank{{Placeholder: "city", Variant: "Москва"}})
	assert.ErrorIs(t, err, answer.ErrInvalid)

	_, err = answer.New([]answer.AnswerBlank{
		{Placeholder: "{{city}}", Variant: "Москва"},
		{Placeholder: "{{city}}", Variant: "Moscow"},
	})
	assert.ErrorIs(t, err, answer.ErrInvalid)
}
