package quiz

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria/random"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuiz(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		err          error
	}{
		{
			name:         "empty title",
			title:        "",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "whitespaces title",
			title:        " ",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "empty sources",
			title:        "title",
			sourceCount:  0,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptySources,
		},
		{
			name:         "negative attempts",
			title:        "title",
			sourceCount:  1,
			attemptLimit: -1,
			timeLimit:    0,
			err:          ErrNegativeAttempts,
		},
		{
			name:         "negative time",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    -1,
			err:          ErrNegativeTime,
		},
		{
			name:         "valid",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			require.NoError(t, err, "failed to create params")

			q, err := NewQuiz(params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, q)

				assert.Equal(t, tt.title, q.Title())
				assert.Len(t, q.Sources(), tt.sourceCount)
				assert.Equal(t, tt.attemptLimit, q.AttemptLimit())
				assert.Equal(t, tt.timeLimit, q.TimeLimit())

				if len(q.Sources()) > 0 {
					// Проверка инкапсуляции: слайсы не должны ссылаться на одну и ту же область памяти
					assert.NotSame(t, &q.Sources()[0], &params.Sources[0], "expected different slices, got equal")
				}
			}
		})
	}
}

func TestIsFiniteAttempts(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true, // В вашем оригинальном коде логика такая (0 = true для IsInfiniteAttempts)
		},
		{
			name:         "finite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 10,
			timeLimit:    0,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			require.NoError(t, err)

			q, err := NewQuiz(params)
			require.NoError(t, err)

			assert.Equal(t, tt.isFinite, q.IsInfiniteAttempts())
		})
	}
}

func TestIsFiniteTime(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true,
		},
		{
			name:         "finite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    10,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			require.NoError(t, err)

			q, err := NewQuiz(params)
			require.NoError(t, err)

			assert.Equal(t, tt.isFinite, q.IsInfiniteTime())
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		err          error
	}{
		{
			name:         "empty title",
			title:        "",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          ErrEmptyTitle,
		},
		{
			name:         "whitespaces title",
			title:        " ",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          ErrEmptyTitle,
		},
		{
			name:         "valid title",
			title:        "new title",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams("title", tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			require.NoError(t, err)

			q, err := NewQuiz(params)
			require.NoError(t, err)

			err = q.Rename(tt.title)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				assert.Equal(t, tt.title, q.Title())
			}
		})
	}
}

func TestAddSource(t *testing.T) {
	t.Run("successive add", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		toAdd, err := createSource(uuid.Nil)
		require.NoError(t, err)

		err = q.AddSource(toAdd)
		
		assert.NoError(t, err)
		assert.Len(t, q.Sources(), 11)
	})

	t.Run("duplicated source", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		toAdd, err := createSource(uuid.Nil)
		require.NoError(t, err)

		err = q.AddSource(toAdd)
		require.NoError(t, err)

		err = q.AddSource(toAdd)
		assert.ErrorIs(t, err, ErrDuplicatedSource)
	})

	t.Run("duplicated bank", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		toAdd, err := createSource(uuid.Nil)
		require.NoError(t, err)

		err = q.AddSource(toAdd)
		require.NoError(t, err)

		dupl, err := createSource(toAdd.BankID())
		require.NoError(t, err)

		err = q.AddSource(dupl)
		assert.ErrorIs(t, err, ErrDuplicatedBank)
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove unexisting elem", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		toRemove, err := createSource(uuid.Nil)
		require.NoError(t, err)

		err = q.RemoveSource(toRemove)
		assert.NoError(t, err) // Удаление несуществующего элемента обычно проходит без ошибок (no-op)
	})

	t.Run("remove existing elem", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		toAdd, err := createSource(uuid.Nil)
		require.NoError(t, err)

		err = q.AddSource(toAdd)
		require.NoError(t, err)
		assert.Len(t, q.Sources(), 11)

		err = q.RemoveSource(toAdd)
		assert.NoError(t, err)
		assert.Len(t, q.Sources(), 10) // Убеждаемся, что элемент реально удалился
	})

	t.Run("remove last elem", func(t *testing.T) {
		params, err := createParams("title", 1, 1, 1)
		require.NoError(t, err)

		q, err := NewQuiz(params)
		require.NoError(t, err)

		err = q.RemoveSource(params.Sources[0])
		assert.ErrorIs(t, err, ErrCannotRemoveLastSource)
	})
}

func TestDelete(t *testing.T) {
	params, err := createParams("title", 1, 1, 1)
	require.NoError(t, err)

	q, err := NewQuiz(params)
	require.NoError(t, err)

	assert.False(t, q.IsDeleted(), "expected quiz be not deleted, got true")

	q.Delete()

	assert.True(t, q.IsDeleted(), "expected quiz be deleted, got false")
}

// === Вспомогательные функции ===

func createParams(title string, sourceCount, attemptLimit, timeLimit int) (Params, error) {
	sources := make([]source.Source, 0, 5)
	for i := 0; i < sourceCount; i++ {
		s, err := createSource(uuid.Nil)
		if err != nil {
			return Params{}, err
		}

		sources = append(sources, s)
	}

	return Params{
		Title:        title,
		Sources:      sources,
		AttemptLimit: attemptLimit,
		TimeLimit:    timeLimit,
	}, nil
}

func createSource(bankID uuid.UUID) (source.Source, error) {
	cParams := random.Params{
		Count: 10,
	}

	r, err := random.NewRandomCriteria(cParams)
	if err != nil {
		return source.Source{}, err
	}

	var params source.Params

	if bankID != uuid.Nil {
		params = source.Params{
			BankID:   bankID,
			Criteria: r,
		}
	} else {
		params = source.Params{
			BankID:   uuid.New(),
			Criteria: r,
		}
	}

	s, err := source.NewSource(params)
	if err != nil {
		return source.Source{}, err
	}

	return s, nil
}
