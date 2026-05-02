package quiz_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("ошибка при пустом списке источников", func(t *testing.T) {
		// Arrange & Act
		q := newQuizBuilder().
			withTitle("Тестовый квиз").
			withSourceList([]source.Source{}). // Пустой список
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, quiz.ErrInvalid)

		// Assert
		assert.Nil(t, q)
	})

	t.Run("ошибка при дублировании bankID", func(t *testing.T) {
		// Arrange & Act
		duplicateID := uuid.New()

		q := newQuizBuilder().
			withTitle("Тестовый квиз").
			withSourceList(mockSourceList(duplicateID, duplicateID, uuid.New())).
			build(t, quiz.ErrInvalid)

		// Assert
		assert.Nil(t, q)
	})

	t.Run("ошибка при превышении лимита источников", func(t *testing.T) {
		// Arrange & Act
		q := newQuizBuilder().
			withTitle("Тестовый квиз").
			withSourceList(mockSourcesWithLength(101)). // Лимит 100
			build(t, quiz.ErrInvalid)

		// Assert
		assert.Nil(t, q)
	})

	t.Run("успешное создание квиза", func(t *testing.T) {
		// Arrange & Act
		q := newQuizBuilder().
			withTitle("Тестовый квиз").
			withSource(mockSource()).
			withMaxAttempts(10).
			withTimeLimit(20).
			build(t, nil)

		// Assert
		assert.NotEqual(t, uuid.Nil, q.ID())
		assert.Equal(t, "Тестовый квиз", q.Title().Value())
		assert.Len(t, q.Sources(), 1)
		assert.Equal(t, 10, q.Attempts().Count())
		assert.Equal(t, 20, q.Time().Seconds())
	})
}

func TestHasInfiniteAttempts(t *testing.T) {
	t.Run("бесконечные попытки", func(t *testing.T) {
		q := newQuizBuilder().withMaxAttempts(0).withSource(mockSource()).build(t, nil)
		assert.True(t, q.HasInfiniteAttempts())
	})

	t.Run("конечные попытки", func(t *testing.T) {
		q := newQuizBuilder().withMaxAttempts(10).withSource(mockSource()).build(t, nil)
		assert.False(t, q.HasInfiniteAttempts())
	})
}

func TestHasInfiniteTime(t *testing.T) {
	t.Run("бесконечное время", func(t *testing.T) {
		q := newQuizBuilder().withTimeLimit(0).withSource(mockSource()).build(t, nil)
		assert.True(t, q.HasInfiniteTime())
	})

	t.Run("конечное время", func(t *testing.T) {
		q := newQuizBuilder().withTimeLimit(20).withSource(mockSource()).build(t, nil)
		assert.False(t, q.HasInfiniteTime())
	})
}

func TestRename(t *testing.T) {
	t.Run("успешное переименование", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().withTitle("Старое название").withSource(mockSource()).build(t, nil)
		newTitle := makeTitle("Новое название")

		// Act
		q.Rename(newTitle)

		// Assert
		assert.Equal(t, "Новое название", q.Title().Value())
	})
}

func TestChangeMaxAttempts(t *testing.T) {
	t.Run("успешное изменение лимита попыток", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().withMaxAttempts(10).withSource(mockSource()).build(t, nil)
		newAttempts := makeAttempts(5)

		// Act
		q.ChangeMaxAttempts(newAttempts)

		// Assert
		assert.Equal(t, 5, q.Attempts().Count())
	})
}

func TestChangeTimeLimit(t *testing.T) {
	t.Run("успешное изменение лимита времени", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().withTimeLimit(20).withSource(mockSource()).build(t, nil)
		newTime := makeTime(300)

		// Act
		q.ChangeTimeLimit(newTime)

		// Assert
		assert.Equal(t, 300, q.Time().Seconds())
	})
}

func TestAddSource(t *testing.T) {
	t.Run("ошибка при добавлении источника с уже существующим банком", func(t *testing.T) {
		// Arrange
		duplicateID := uuid.New()
		q := newQuizBuilder().
			withSourceList(mockSourceList(duplicateID, uuid.New())).
			build(t, nil)

		sourceToAdd := mockSourceList(duplicateID)[0]

		// Act
		err := q.AddSource(sourceToAdd)

		// Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Contains(t, err.Error(), "источник с таким банком уже указан")
		assert.Len(t, q.Sources(), 2) // Состояние не изменилось
	})

	t.Run("ошибка при превышении максимального количества источников", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().
			withSourceList(mockSourcesWithLength(100)).
			build(t, nil)

		// Act
		err := q.AddSource(mockSource())

		// Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Contains(t, err.Error(), "максимальное количество")
		assert.Len(t, q.Sources(), 100)
	})

	t.Run("успешное добавление источника", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().
			withSourceList(mockSourcesWithLength(2)).
			build(t, nil)

		// Act
		err := q.AddSource(mockSource())

		// Assert
		assert.NoError(t, err)
		assert.Len(t, q.Sources(), 3)
	})
}

func TestRemoveSource(t *testing.T) {
	t.Run("ошибка при удалении последнего источника", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().
			withSourceList(mockSourcesWithLength(1)).
			build(t, nil)

		sourceToRemove := q.Sources()[0]

		// Act
		err := q.RemoveSource(sourceToRemove)

		// Assert
		assert.ErrorIs(t, err, quiz.ErrInvalid)
		assert.Contains(t, err.Error(), "нельзя удалить последний источник")
		assert.Len(t, q.Sources(), 1)
	})

	t.Run("успешное выполнение (игнорирование), если источника нет в квизе", func(t *testing.T) {
		// Arrange
		q := newQuizBuilder().
			withSourceList(mockSourcesWithLength(5)).
			build(t, nil)

		nonExistingSource := mockSource()

		// Act
		err := q.RemoveSource(nonExistingSource)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, q.Sources(), 5)
	})

	t.Run("успешное удаление источника", func(t *testing.T) {
		// Arrange
		target := mockSourceList(uuid.New())[0]

		q := newQuizBuilder().
			withSourceList(mockSourceList(target.BankID(), uuid.New(), uuid.New())).
			build(t, nil)

		// Act
		err := q.RemoveSource(target)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, q.Sources(), 2)

		// Убеждаемся, что именно нужный источник удален
		for _, s := range q.Sources() {
			assert.NotEqual(t, target.BankID(), s.BankID())
		}
	})
}
