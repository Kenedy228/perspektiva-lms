package attempt

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		name         string
		questionType question.Type
		snapshot     []byte
		err          error
	}{
		{
			name:         "invalid questionType",
			questionType: question.Type("invalid"),
			snapshot:     []byte("payload"),
			err:          ErrInvalidQuestionType,
		},
		{
			name:         "empty snapshot",
			questionType: question.TypeMatching,
			snapshot:     []byte(""),
			err:          ErrEmptySnapshot,
		},
		{
			name:         "valid data",
			questionType: question.TypeMatching,
			snapshot:     []byte("payload"),
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewItem(tt.questionType, tt.snapshot)

			assert.ErrorIs(t, err, tt.err)
			require.NotNil(t, item)

			if err == nil {
				assert.NotEqual(t, item.ID(), uuid.Nil)
				assert.Equal(t, item.Snapshot(), tt.snapshot)
				assert.Equal(t, item.QuestionType(), tt.questionType.String())
				assert.Nil(t, item.StudentAnswer())
				assert.Nil(t, item.Score())
				assert.False(t, item.IsAnswered())
				assert.NotSame(t, &item.Snapshot()[0], &tt.snapshot[0])
			}
		})
	}
}
