package attempt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	t.Run("nil question", func(t *testing.T) {
		_, err := newItem(nil)

		assert.ErrorIs(t, err, ErrInvalidItem)
	})

	t.Run("valid question", func(t *testing.T) {
		mockQ := new(mockQuestion)
		mockQ.On("Clone").Return(mockQ)

		item, err := newItem(mockQ)

		assert.Nil(t, err)
		assert.Equal(t, item.Score(), 0)
		assert.Nil(t, item.Answer())
		assert.False(t, item.HasAnswer())
	})
}

func TestChangeAnswer(t *testing.T) {
	t.Run("nil answer", func(t *testing.T) {
		mockQ := new(mockQuestion)
		mockQ.On("Clone").Return(mockQ)

		item, err := newItem(mockQ)
		assert.Nil(t, err)

		assert.Nil(t, item.Answer())
		assert.False(t, item.HasAnswer())
		assert.Equal(t, item.calculateScore(), 0)
	})

	t.Run("non-nil answer", func(t *testing.T) {
		mockQ := new(mockQuestion)
		mockQ.On("Clone").Return(mockQ)

		item, err := newItem(mockQ)
		assert.Nil(t, err)

		oldChangedAt := item.ChangedAt()

		answer := new(mockAnswer)
		answer.On("Clone").Return(answer)
		item.changeAnswer(answer)

		assert.NotNil(t, item.Answer())
		assert.True(t, item.HasAnswer())
		assert.True(t, item.ChangedAt().After(oldChangedAt))
	})
}

func TestCalculateScore(t *testing.T) {
	t.Run("correct answer", func(t *testing.T) {
		answer := new(mockAnswer)
		answer.On("Clone").Return(answer)

		mockQ := new(mockQuestion)
		mockQ.On("Clone").Return(mockQ)
		mockQ.On("CheckAnswer", answer).Return(true)

		item, err := newItem(mockQ)
		assert.Nil(t, err)

		item.changeAnswer(answer)

		assert.True(t, item.HasAnswer())

		score := item.calculateScore()
		assert.Equal(t, score, 1)
		assert.Equal(t, item.Score(), 1)
	})

	t.Run("incorrect answer", func(t *testing.T) {
		answer := new(mockAnswer)
		answer.On("Clone").Return(answer)

		mockQ := new(mockQuestion)
		mockQ.On("Clone").Return(mockQ)
		mockQ.On("CheckAnswer", answer).Return(true)

		item, err := newItem(mockQ)
		assert.Nil(t, err)

		item.changeAnswer(answer)

		assert.True(t, item.HasAnswer())
		score := item.calculateScore()
		assert.Equal(t, score, 1)
		assert.Equal(t, item.Score(), 1)
	})
}

func TestClone(t *testing.T) {
	answer := new(mockAnswer)
	answer.On("Clone").Return(answer)

	mockQ := new(mockQuestion)
	mockQ.On("Clone").Return(mockQ)

	item, err := newItem(mockQ)
	item.changeAnswer(answer)

	assert.Nil(t, err)

	cItem := item.Clone()

	assert.NotSame(t, &item, &cItem)
}
