package score_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s, err := score.New(0.5)
	require.NoError(t, err)
	assert.Equal(t, 0.5, s.Value())

	_, err = score.New(-0.1)
	assert.ErrorIs(t, err, score.ErrInvalid)

	_, err = score.New(1.1)
	assert.ErrorIs(t, err, score.ErrInvalid)
}
