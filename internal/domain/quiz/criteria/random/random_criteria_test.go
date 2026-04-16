package random

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomCriteria(t *testing.T) {
	tests := []struct {
		name  string
		count int
		err   error
	}{
		{
			name:  "negative count",
			count: -1,
			err:   ErrInvalidCount,
		},
		{
			name:  "zero count",
			count: 0,
			err:   ErrInvalidCount,
		},
		{
			name:  "valid count",
			count: 1,
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := Params{
				Count: tt.count,
			}
			
			c, err := NewRandomCriteria(params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				// Убеждаемся, что при отсутствии ошибки объект точно создался
				require.NotNil(t, c)
				
				assert.Equal(t, criteria.CriteriaTypeRandom, c.Type())
				assert.Equal(t, tt.count, c.QuestionCount())
			}
		})
	}
}
