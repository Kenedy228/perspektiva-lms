package bank_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/bank"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("empty title should return err", func(t *testing.T) {
		//Arrange-Assert
		newBankBuilder().build(t, bank.ErrInvalid)
		newBankBuilder().withTitle(" ").build(t, bank.ErrInvalid)
	})

	t.Run("success", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)

		//Assert
		assert.NotEqual(t, b.ID(), uuid.Nil)
		assert.Equal(t, b.Title(), "title")
		assert.Equal(t, len(b.Questions()), 0)
		assert.False(t, b.CreatedAt().IsZero())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
		assert.True(t, b.DeletedAt().IsZero())
		assert.False(t, b.IsDeleted())
	})
}

func TestRename(t *testing.T) {
	t.Run("empty title return err and not update", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.Rename("")

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.Rename("new title")

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, b.Title(), "new title")
		assert.True(t, oldUpdatedAt.Before(b.UpdatedAt()))
		assert.True(t, b.UpdatedAt().After(b.CreatedAt()))
	})
}

func TestAddQuestions(t *testing.T) {
	t.Run("adding zero questions should not return err and update state", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.AddQuestions()

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, len(b.Questions()), 0)
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("max len exceeded return error", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e6)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("contains nil uuid return error", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)
		ids = append(ids, uuid.Nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("adding questions contains duplicates", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		ids := makeIDs(1e3)
		ids = append(ids, ids[0])
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.AddQuestions(ids...)

		//Assert
		assert.ErrorIs(t, err, bank.ErrInvalid)
		assert.Equal(t, len(b.Questions()), 0)
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("bank questions contains duplicates of adding", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			ids := makeIDs(1e3)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			b.AddQuestions(ids...)
			err := b.AddQuestions(ids...)

			//Assert
			assert.ErrorIs(t, err, bank.ErrInvalid)
			assert.Equal(t, len(b.Questions()), len(ids))
			assert.True(t, oldUpdatedAt.Before(b.UpdatedAt()))
			assert.True(t, b.UpdatedAt().After(b.CreatedAt()))
		})
	})

	t.Run("valid", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			ids := makeIDs(1e3)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := b.AddQuestions(ids...)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(b.Questions()), len(ids))
			assert.True(t, oldUpdatedAt.Before(b.UpdatedAt()))
			assert.True(t, b.UpdatedAt().After(b.CreatedAt()))
		})
	})
}

func TestRemoveQuestions(t *testing.T) {
	t.Run("should not change state if len of questions to delete is zero", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.RemoveQuestions()

		//Assert
		assert.NoError(t, err)
		assert.Empty(t, b.Questions())
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("should not change state if not delete any question", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()
		ids := makeIDs(100)

		//Act
		err := b.RemoveQuestions(ids...)

		//Assert
		assert.NoError(t, err)
		assert.Empty(t, b.Questions())
		assert.Equal(t, oldUpdatedAt, b.UpdatedAt())
		assert.Equal(t, b.CreatedAt(), b.UpdatedAt())
	})

	t.Run("should change state if any question was removed", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()
			ids := makeIDs(100)

			//Act
			time.Sleep(time.Second * 100)
			b.AddQuestions(ids...)
			err := b.RemoveQuestions(ids...)

			//Assert
			assert.NoError(t, err)
			assert.Empty(t, b.Questions())
			assert.True(t, b.UpdatedAt().After(oldUpdatedAt))
			assert.True(t, b.CreatedAt().Before(b.UpdatedAt()))
		})
	})
}

func TestClearQuestions(t *testing.T) {
	t.Run("clear questions should remove all elements", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()
			ids := makeIDs(100)

			//Act
			time.Sleep(time.Second * 100)
			b.AddQuestions(ids...)
			err := b.ClearQuestions()

			//Assert
			assert.NoError(t, err)
			assert.Empty(t, b.Questions())
			assert.True(t, b.UpdatedAt().After(oldUpdatedAt))
			assert.True(t, b.CreatedAt().Before(b.UpdatedAt()))
		})
	})

	t.Run("clear questions should work on empty questions", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second * 100)
			err := b.ClearQuestions()

			//Assert
			assert.NoError(t, err)
			assert.Empty(t, b.Questions())
			assert.True(t, b.UpdatedAt().After(oldUpdatedAt))
			assert.True(t, b.CreatedAt().Before(b.UpdatedAt()))
		})
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete non-deleted bank and permit any operation (besides getters and Restore)", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)
			oldUpdatedAt := b.UpdatedAt()

			//Act
			time.Sleep(time.Second * 100)
			err := b.Delete()

			//Assert
			assert.NoError(t, err)
			assert.True(t, b.IsDeleted())
			assert.False(t, b.DeletedAt().IsZero())
			assert.True(t, b.UpdatedAt().After(oldUpdatedAt))
			assert.True(t, b.CreatedAt().Before(b.UpdatedAt()))
			assert.ErrorIs(t, b.AddQuestions(), bank.ErrDeleted)
			assert.ErrorIs(t, b.RemoveQuestions(), bank.ErrDeleted)
			assert.ErrorIs(t, b.Delete(), bank.ErrDeleted)
			assert.ErrorIs(t, b.Rename(""), bank.ErrDeleted)
			assert.ErrorIs(t, b.ClearQuestions(), bank.ErrDeleted)
		})
	})

	t.Run("should not delete deleted bank", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			b := newBankBuilder().withTitle("title").build(t, nil)

			//Act
			time.Sleep(time.Second * 100)
			b.Delete()
			oldUpdatedAt := b.UpdatedAt()
			err := b.Delete()

			//Assert
			assert.ErrorIs(t, err, bank.ErrDeleted)
			assert.True(t, b.IsDeleted())
			assert.False(t, b.DeletedAt().IsZero())
			assert.Equal(t, b.UpdatedAt(), oldUpdatedAt)
			assert.ErrorIs(t, b.AddQuestions(), bank.ErrDeleted)
			assert.ErrorIs(t, b.RemoveQuestions(), bank.ErrDeleted)
			assert.ErrorIs(t, b.Delete(), bank.ErrDeleted)
			assert.ErrorIs(t, b.Rename(""), bank.ErrDeleted)
		})
	})
}

func TestRestore(t *testing.T) {
	t.Run("should return err if bank is not deleted", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		err := b.Restore()

		//Assert
		assert.ErrorIs(t, err, bank.ErrNotDeleted)
		assert.False(t, b.IsDeleted())
		assert.Equal(t, b.UpdatedAt(), oldUpdatedAt)
	})

	t.Run("should restore bank if it was deleted", func(t *testing.T) {
		//Arrange
		b := newBankBuilder().withTitle("title").build(t, nil)
		oldUpdatedAt := b.UpdatedAt()

		//Act
		b.Delete()
		err := b.Restore()

		//Assert
		assert.NoError(t, err)
		assert.False(t, b.IsDeleted())
		assert.True(t, b.UpdatedAt().After(oldUpdatedAt))
	})
}
