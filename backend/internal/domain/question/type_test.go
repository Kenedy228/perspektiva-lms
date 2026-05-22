package question_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType(t *testing.T) {
	got, err := question.ParseType("selectable")
	require.NoError(t, err)
	assert.Equal(t, question.TypeSelectable, got)
	assert.True(t, got.IsValid())

	_, err = question.ParseType("unknown")
	assert.Error(t, err)
	assert.False(t, question.Type("unknown").IsValid())
}
