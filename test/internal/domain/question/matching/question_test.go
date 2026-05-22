//go:build legacy
// +build legacy

package matching_test

import (
	"testing"

	matching2 "gitflic.ru/lms/backend/internal/domain/question/matching"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		count int
	}{
		{
			name:  "количество пар удовлетворяет инвариантам",
			count: 10,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := newQuestionBuilder().
				withPairs(tt.count).
				withTitle("title").
				build()

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, q.Title().Value(), "title")
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name    string
		count   int
		wantErr error
	}{
		{
			name:    "количество пар меньше необходимого",
			count:   1,
			wantErr: matching2.ErrInvalid,
		},
		{
			name:    "количество пар больше необходимого",
			count:   30,
			wantErr: matching2.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := newQuestionBuilder().
				withPairs(tt.count).
				withTitle("title").
				build()

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestMatchingQuestion_UpdatePairs(t *testing.T) {
	//Arrange
	q, err := newQuestionBuilder().
		withPairs(10).
		withTitle("title").
		build()
	require.NoError(t, err)

	//Act
	pairs := makePairs(10)
	q.ChangePairs(pairs)

	//Assert
	assert.Equal(t, q.Pairs(), pairs)
}

func TestClone(t *testing.T) {
	//Arrange
	q, err := newQuestionBuilder().
		withPairs(10).
		withTitle("title").
		build()
	require.NoError(t, err)
	clone, ok := q.Clone().(*matching2.Question)
	require.True(t, ok)

	//Assert
	assert.Equal(t, clone.ID(), q.ID())
	assert.Equal(t, clone.Title(), q.Title())
	assert.Equal(t, clone.Instruction(), q.Instruction())
	assert.Equal(t, clone.Pairs(), q.Pairs())
	assert.NotSame(t, &clone.Pairs()[0], &q.Pairs()[0])
}
