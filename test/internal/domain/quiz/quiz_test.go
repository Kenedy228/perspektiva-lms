package quiz_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/quiz"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewQuiz(t *testing.T) {
	t.Run("with empty title returns error", func(t *testing.T) {
		//Arrange-Assert
		newQuizBuilder().withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)

		newQuizBuilder().withSource(mockSource()).
			withTitle(" ").
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("with empty sources returns error", func(t *testing.T) {
		//Arrange-Assert
		newQuizBuilder().withTitle("title").
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("with duplicated bankID returns error", func(t *testing.T) {
		//Arrange-Assert
		duplicateID := uuid.New()

		newQuizBuilder().withTitle("title").
			withSourceList(mockSourceList(duplicateID, duplicateID, uuid.New())).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("with maxSources sources returns error", func(t *testing.T) {
		//Arrange-Assert
		newQuizBuilder().withTitle("title").
			withSourceList(mockSourcesWithLength(1e4)).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("with negative maxAttempts return error", func(t *testing.T) {
		//Arrange-Assert
		newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(-1).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("with maxAttempts violation return error", func(t *testing.T) {
		//Arrange-Assert
		newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(1e5).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Assert
		assert.NotEqual(t, q.ID(), uuid.Nil)
		assert.Equal(t, q.Title(), "title")
		assert.Equal(t, len(q.Sources()), 1)
		assert.Equal(t, q.MaxAttempts(), 10)
		assert.Equal(t, q.TimeLimit().Value(), 20)
		assert.False(t, q.CreatedAt().IsZero())
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})
}

func TestHasInfiniteAttempts(t *testing.T) {
	t.Run("infinite", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(0).
			withTimeLimit(20).
			build(t, nil)

		//Assert
		assert.True(t, q.HasInfiniteAttempts())
	})

	t.Run("finite", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Assert
		assert.False(t, q.HasInfiniteAttempts())
	})
}

func TestHasInfiniteTime(t *testing.T) {
	t.Run("infinite", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(0).
			build(t, nil)

		//Assert
		assert.True(t, q.HasInfiniteTime())
	})

	t.Run("finite", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Assert
		assert.False(t, q.HasInfiniteTime())
	})
}

func TestRename(t *testing.T) {
	t.Run("invalid title return error", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(0).
			build(t, nil)

		//Act
		err := q.Rename("")

		//Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Equal(t, q.Title(), "title")
		assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
	})

	t.Run("valid title success", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := newQuizBuilder().withTitle("title").
				withSource(mockSource()).
				withMaxAttempts(10).
				withTimeLimit(0).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			err := q.Rename("new title")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, q.Title(), "new title")
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()), "title")
		})
	})
}

func TestAddSource(t *testing.T) {
	t.Run("add source with duplicate bank id return error", func(t *testing.T) {
		//Arrange
		duplicateID := uuid.New()

		q := newQuizBuilder().withTitle("title").
			withSourceList(mockSourceList(duplicateID, uuid.New())).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.AddSource(mockSourceList(duplicateID)[0])

		//Assert
		assert.ErrorIs(t, err, quiz.ErrDuplicateBankID)
		assert.Equal(t, len(q.Sources()), 2)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("add source with sources limit return error", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSourceList(mockSourcesWithLength(100)).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.AddSource(mockSource())

		//Assert
		assert.ErrorIs(t, err, quiz.ErrSourceSizeExceeded)
		assert.Equal(t, len(q.Sources()), 100)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("valid", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := newQuizBuilder().withTitle("title").
				withSourceList(mockSourcesWithLength(20)).
				withMaxAttempts(10).
				withTimeLimit(20).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			err := q.AddSource(mockSource())

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(q.Sources()), 21)
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
		})
	})
}

func TestRemoveSource(t *testing.T) {
	t.Run("remove last source return error", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSourceList(mockSourcesWithLength(1)).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.RemoveSource(mockSource())

		//Assert
		assert.ErrorIs(t, err, quiz.ErrCannotRemoveLastSource)
		assert.Equal(t, len(q.Sources()), 1)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("remove source with non-existing resource doesn't return err", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSourceList(mockSourcesWithLength(10)).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.RemoveSource(mockSource())

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, len(q.Sources()), 10)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("success", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			target := mockSourceList(uuid.New())[0]

			q := newQuizBuilder().withTitle("title").
				withSourceList(mockSourceList(target.BankID(), uuid.New())).
				withMaxAttempts(10).
				withTimeLimit(20).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			err := q.RemoveSource(target)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(q.Sources()), 1)
			assert.NotContains(t, q.Sources(), target)
			assert.NotEqual(t, q.Sources()[0].BankID(), uuid.Nil)
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
		})
	})
}

func TestChangeMaxAttempts(t *testing.T) {
	t.Run("with negative maxAttempts return error", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.ChangeMaxAttempts(-1)

		//Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Equal(t, q.MaxAttempts(), 10)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("with maxAttempts violation return error", func(t *testing.T) {
		//Arrange
		q := newQuizBuilder().withTitle("title").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		//Act
		err := q.ChangeMaxAttempts(1e5)

		//Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Equal(t, q.MaxAttempts(), 10)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})

	t.Run("success", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := newQuizBuilder().withTitle("title").
				withSource(mockSource()).
				withMaxAttempts(10).
				withTimeLimit(20).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			err := q.ChangeMaxAttempts(20)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, q.MaxAttempts(), 20)
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
		})
	})
}

func TestChangeTimeLimit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := newQuizBuilder().withTitle("title").
				withSource(mockSource()).
				withMaxAttempts(10).
				withTimeLimit(20).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			q.ChangeTimeLimit(makeLimit(300))

			//Assert
			assert.Equal(t, q.TimeLimit().Value(), 300)
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
		})
	})
}
