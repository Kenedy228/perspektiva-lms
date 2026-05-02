package bank_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/bank"
	"gitflic.ru/lms/internal/domain/bank/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("пустое название должно вернуть ошибку", func(t *testing.T) {
		//Arrange-Assert
		newBankBuilder().build(t, title.ErrInvalid)
		newBankBuilder().withTitle(" ").build(t, title.ErrInvalid)
	})

	t.Run("успех", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Assert
		assert.NotEqual(t, b.ID(), uuid.Nil)
		assert.Equal(t, b.Title().Value(), "title")
		assert.Equal(t, len(b.Questions()), 0)
	})
}

func TestRename(t *testing.T) {
	t.Run("валидно", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Act
		newTitle, err := title.New("new title")
		assert.NoError(t, err)
		b.Rename(newTitle)

		//Assert
		assert.Equal(t, b.Title().Value(), "new title")
	})
}

func TestAddQuestions(t *testing.T) {
	t.Run("добавление нуля вопросов не должно вернуть ошибку", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Act
		err := b.AddQuestions()

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, len(b.Questions()), 0)
	})

	t.Run("превышение максимальной длины возвращает ошибку", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e6)

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
	})

	t.Run("содержит nil uuid возвращает ошибку", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)
		ids = append(ids, uuid.Nil)

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
	})

	t.Run("добавляемые вопросы содержат дубликаты", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)
		ids = append(ids, ids[0])

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
	})

	t.Run("вопросы банка содержат дубликаты добавляемых", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)

		//Act
		b.AddQuestions(ids...)
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), len(ids))
	})

	t.Run("валидно", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, len(b.Questions()), len(ids))
	})
}

func TestRemoveQuestions(t *testing.T) {
	t.Run("не должно менять состояние если длина удаляемых вопросов ноль", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Act
		b.RemoveQuestions()

		//Assert
		assert.Empty(t, b.Questions())
	})

	t.Run("не должно менять состояние если не удален ни один вопрос", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(100)

		//Act
		b.RemoveQuestions(ids...)

		//Assert
		assert.Empty(t, b.Questions())
	})

	t.Run("должен удалить существующие вопросы", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(100)

		//Act
		b.AddQuestions(ids...)
		b.RemoveQuestions(ids...)

		//Assert
		assert.Empty(t, b.Questions())
	})

	t.Run("должен удалить только существующие вопросы", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(100)
		nonExisting := makeIDs(50)

		//Act
		b.AddQuestions(ids...)
		b.RemoveQuestions(append(ids[:10], nonExisting...)...)

		//Assert
		assert.Equal(t, len(b.Questions()), 90)
	})
}

func TestClearQuestions(t *testing.T) {
	t.Run("очистка вопросов должна удалить все элементы", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(100)

		//Act
		b.AddQuestions(ids...)
		b.ClearQuestions()

		//Assert
		assert.Empty(t, b.Questions())
	})

	t.Run("очистка вопросов должна работать на пустых вопросах", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Act
		b.ClearQuestions()

		//Assert
		assert.Empty(t, b.Questions())
	})
}
