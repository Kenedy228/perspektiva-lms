package manual

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManualCriteria(t *testing.T) {
	tests := []struct {
		name string
		ids  []uuid.UUID
		err  error
	}{
		{
			name: "empty ids",
			ids:  []uuid.UUID{},
			err:  ErrEmptyQuestions,
		},
		{
			name: "nil id",
			ids:  []uuid.UUID{uuid.Nil},
			err:  ErrNilQuestion,
		},
		{
			name: "valid ids",
			ids:  []uuid.UUID{uuid.New(), uuid.New()},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := Params{
				QuestionIDs: tt.ids,
			}

			c, err := NewManualCriteria(params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, c)

				assert.Equal(t, criteria.CriteriaTypeManual, c.Type())
				assert.Equal(t, len(tt.ids), c.QuestionCount())

				if c.QuestionCount() > 0 {
					ids := c.QuestionIDs()
					
					// Проверяем, что содержимое слайсов совпадает
					assert.Equal(t, tt.ids, ids)
					
					// Проверяем, что это копия (адреса первых элементов отличаются)
					assert.NotSame(t, &tt.ids[0], &ids[0], "expected to find copy of slice, got the same slice")
				}
			}
		})
	}
}
