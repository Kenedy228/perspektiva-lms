package base_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/question/base"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create with nil image", func(t *testing.T) {
		//Arrange
		base := newBaseBuilder().withText("text").build(t, nil)

		//Assert
		assert.NotEqual(t, base.ID(), uuid.Nil)
		assert.Equal(t, base.Text(), "text")
		assert.Equal(t, base.ImageID(), uuid.Nil)
		assert.False(t, base.HasImage())
		assert.Equal(t, base.UpdatedAt(), base.CreatedAt())
	})

	t.Run("should create", func(t *testing.T) {
		//Arrange
		base := newBaseBuilder().withText("text").withImage(uuid.New()).build(t, nil)

		//Assert
		assert.NotEqual(t, base.ID(), uuid.Nil)
		assert.Equal(t, base.Text(), "text")
		assert.NotEqual(t, base.ImageID(), uuid.Nil)
		assert.True(t, base.HasImage())
		assert.Equal(t, base.UpdatedAt(), base.CreatedAt())
	})

	t.Run("should return err if text empty or has only whitespaces", func(t *testing.T) {
		//Arrange-Assert
		newBaseBuilder().withImage(uuid.New()).build(t, base.ErrInvalid)
		newBaseBuilder().withText(" ").withImage(uuid.New()).build(t, base.ErrInvalid)
	})
}

func TestRemoveImage(t *testing.T) {
	t.Run("should remove image if image is set and update updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").withImage(uuid.New()).build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.RemoveImage()

			//Assert
			assert.Equal(t, base.ImageID(), uuid.Nil)
			assert.False(t, base.HasImage())
			assert.True(t, oldUpdatedAt.Before(base.UpdatedAt()))
		})
	})

	t.Run("should not remove image if image is not set and not update updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.RemoveImage()

			//Assert
			assert.Equal(t, base.ImageID(), uuid.Nil)
			assert.False(t, base.HasImage())
			assert.Equal(t, oldUpdatedAt, base.UpdatedAt())
		})
	})
}

func TestUpdateImage(t *testing.T) {
	t.Run("should update image if set and change updated at", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			id := uuid.New()
			base := newBaseBuilder().withText("text").withImage(id).build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.UpdateImage(uuid.New())

			//Assert
			assert.NotEqual(t, base.ImageID(), id)
			assert.True(t, base.HasImage())
			assert.True(t, oldUpdatedAt.Before(base.UpdatedAt()))
		})
	})

	t.Run("should update image if not set and change updated at", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.UpdateImage(uuid.New())

			//Assert
			assert.NotEqual(t, base.ImageID(), uuid.Nil)
			assert.True(t, base.HasImage())
			assert.True(t, oldUpdatedAt.Before(base.UpdatedAt()))
		})
	})

	t.Run("should remove image if passed id is nil and change updated at", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").withImage(uuid.New()).build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.UpdateImage(uuid.Nil)

			//Assert
			assert.Equal(t, base.ImageID(), uuid.Nil)
			assert.False(t, base.HasImage())
			assert.True(t, oldUpdatedAt.Before(base.UpdatedAt()))
		})
	})

	t.Run("should not remove nil-image if passed id is nil and not change updated at", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			base.UpdateImage(uuid.Nil)

			//Assert
			assert.Equal(t, base.ImageID(), uuid.Nil)
			assert.False(t, base.HasImage())
			assert.Equal(t, oldUpdatedAt, base.UpdatedAt())
		})
	})
}

func TestUpdateText(t *testing.T) {
	t.Run("should update text", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			base := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := base.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			err := base.UpdateText("new text")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, base.Text(), "new text")
			assert.True(t, oldUpdatedAt.Before(base.UpdatedAt()))
		})
	})

	t.Run("should return err if text is empty, not change text and updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			err := b.UpdateText("")

			//Assert
			assert.ErrorIs(t, err, base.ErrInvalid)
			assert.Equal(t, b.Text(), "text")
			assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		})
	})

	t.Run("should return err if text is whitespaces, not change text and updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBaseBuilder().withText("text").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second)
			err := b.UpdateText(" ")

			//Assert
			assert.ErrorIs(t, err, base.ErrInvalid)
			assert.Equal(t, b.Text(), "text")
			assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		})
	})
}

func TestClone(t *testing.T) {
	//Arrange
	b := newBaseBuilder().withText("text").build(t, nil)
	clone := b.Clone()

	//Assert
	assert.Equal(t, b.ID(), clone.ID())
	assert.Equal(t, b.Text(), clone.Text())
	assert.Equal(t, b.ImageID(), clone.ImageID())
	assert.Equal(t, b.CreatedAt(), clone.CreatedAt())
	assert.Equal(t, b.UpdatedAt(), clone.UpdatedAt())
}
