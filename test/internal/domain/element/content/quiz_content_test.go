package content_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/element/content"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuizContent(t *testing.T) {
	tests := []struct {
		name    string
		quizID  uuid.UUID
		wantErr error
	}{
		{
			name:    "успешное создание quiz content",
			quizID:  uuid.New(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrrange
			got, err := content.NewQuizContent(tt.quizID)

			//Assert
			assert.NoError(t, err)
			assert.NotEqual(t, got.QuizID(), uuid.Nil)
		})
	}
}

func TestNewQuizContent_Fail(t *testing.T) {
	tests := []struct {
		name    string
		quizID  uuid.UUID
		wantErr error
	}{
		{
			name:    "ошибка если quiz id nil",
			quizID:  uuid.Nil,
			wantErr: content.ErrEmptyQuizID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := content.NewQuizContent(tt.quizID)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestQuizContent_Clone(t *testing.T) {
	//Arrange
	original, err := content.NewQuizContent(uuid.New())
	require.NoError(t, err)
	cloned, ok := original.Clone().(content.QuizContent)
	require.True(t, ok)

	//Assert
	assert.Equal(t, original, cloned)
	assert.Equal(t, original.QuizID(), cloned.QuizID())
	assert.Equal(t, original.Type(), cloned.Type())
	assert.Equal(t, original.IsInteractive(), cloned.IsInteractive())
}
