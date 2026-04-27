package element_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/element"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewElement_Success(t *testing.T) {
	//Arrange
	el := newElementBuilder().
		withTitle("  Основы Go  ").
		withContent().
		build(t, nil)

	//Assert
	assert.NotEqual(t, uuid.Nil, el.ID())
	assert.Equal(t, "Основы Go", el.Title(), "заголовок должен быть нормализован")
	assert.False(t, el.CreatedAt().IsZero())
	assert.False(t, el.UpdatedAt().IsZero())
	assert.Equal(t, el.CreatedAt(), el.UpdatedAt())
	assert.NotNil(t, el.Content())
}

func TestNewElement_ValidationErrors(t *testing.T) {
	t.Run("пустой заголовок", func(t *testing.T) {
		//Arrange-Assert
		newElementBuilder().withContent().build(t, element.ErrInvalid)
	})

	t.Run("nil контент", func(t *testing.T) {
		//Arrange-Assert
		newElementBuilder().withTitle("title").build(t, element.ErrInvalid)
	})
}

func TestElement_Mutations(t *testing.T) {
	t.Run("ChangeTitle", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			el := newElementBuilder().withTitle("old title").withContent().build(t, nil)
			oldUpdatedAt := el.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := el.ChangeTitle("new title")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, el.Title(), "new title")
			assert.True(t, el.UpdatedAt().After(oldUpdatedAt))
		})
	})

	t.Run("ChangeTitle with error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			el := newElementBuilder().withTitle("old title").withContent().build(t, nil)
			oldUpdatedAt := el.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := el.ChangeTitle("")

			//Assert
			assert.ErrorIs(t, err, element.ErrInvalid)
			assert.Equal(t, el.Title(), "old title")
			assert.Equal(t, oldUpdatedAt, el.UpdatedAt())
		})
	})

	t.Run("ChangeContent", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			el := newElementBuilder().withTitle("title").withContent().build(t, nil)
			oldUpdatedAt := el.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := el.ChangeContent(&contentFixture{})

			//Assert
			assert.NoError(t, err)
			assert.True(t, el.UpdatedAt().After(oldUpdatedAt))
		})
	})

	t.Run("ChangeContent with error", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			el := newElementBuilder().withTitle("title").withContent().build(t, nil)
			oldUpdatedAt := el.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := el.ChangeContent(nil)

			//Assert
			assert.ErrorIs(t, err, element.ErrInvalid)
			assert.Equal(t, oldUpdatedAt, el.UpdatedAt())
		})
	})
}

func TestElement_CopyAndClone(t *testing.T) {
	t.Run("Copy создает новую сущность", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			original := newElementBuilder().withTitle("title").withContent().build(t, nil)

			//Act
			time.Sleep(time.Second * 10)
			copied, err := original.Copy()

			//Assert
			assert.NoError(t, err)
			assert.NotEqual(t, original.ID(), copied.ID(), "ID должны отличаться")
			assert.Equal(t, original.Title(), copied.Title(), "заголовки совпадают")
			assert.True(t, copied.CreatedAt().After(original.CreatedAt()), "создано позже")

			err = copied.ChangeTitle("Измененная копия")

			assert.NoError(t, err)
			assert.Equal(t, "title", original.Title())
		})
	})

	t.Run("Clone создает техническую копию", func(t *testing.T) {
		//Arrange
		original := newElementBuilder().withTitle("title").withContent().build(t, nil)
		cloned := original.Clone()

		//Assert
		assert.Equal(t, original.ID(), cloned.ID(), "ID должны совпадать")
		assert.Equal(t, original.Title(), cloned.Title())
		assert.Equal(t, original.CreatedAt(), cloned.CreatedAt())

		// Контент клонируется (deep copy)
		assert.Equal(t, original.Content(), cloned.Content())
	})
}
