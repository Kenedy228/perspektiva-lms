package person_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/person/name"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPerson_Success(t *testing.T) {
	t.Run("создание человека без профиля", func(t *testing.T) {
		//Arrange
		p := newPersonBuilder().
			withName().
			build()

		//Act
		_, ok := p.Profile()

		//Assert
		assert.NotEqual(t, uuid.Nil, p.ID())
		assert.Equal(t, nameFixture(), p.Name())
		assert.False(t, p.CreatedAt().IsZero())
		assert.False(t, p.UpdatedAt().IsZero())
		assert.Equal(t, p.CreatedAt(), p.UpdatedAt())
		assert.False(t, p.HasProfile())
		assert.False(t, ok)
	})

	t.Run("создание человека с профилем", func(t *testing.T) {
		//Arrange
		p := newPersonBuilder().
			withName().
			withProfile().
			build()

		//Act
		_, ok := p.Profile()

		//Assert
		assert.True(t, p.HasProfile())
		assert.True(t, ok)
	})
}

func TestPerson_Mutations(t *testing.T) {
	t.Run("AttachProfile добавляет профиль и обновляет updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			p := newPersonBuilder().
				withName().
				build()

			//Act
			initialUpdatedAt := p.UpdatedAt()
			time.Sleep(10 * time.Second)
			newProf := profileFixture()
			p.AttachProfile(newProf)
			_, ok := p.Profile()

			//Assert
			assert.True(t, p.HasProfile())
			assert.True(t, p.UpdatedAt().After(initialUpdatedAt))
			assert.True(t, ok)
		})
	})

	t.Run("DetachProfile удаляет профиль, если он уже установлен, и обновляет updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			p := newPersonBuilder().
				withName().
				build()

			//Act
			newProf := profileFixture()
			p.AttachProfile(newProf)

			attachedUpdatedAt := p.UpdatedAt()
			time.Sleep(10 * time.Second)
			p.DetachProfile()
			_, ok := p.Profile()

			//Assert
			assert.False(t, p.HasProfile())
			assert.True(t, p.UpdatedAt().After(attachedUpdatedAt))
			assert.False(t, ok)
		})
	})

	t.Run("DetachProfile не обновляет updatedAt, если профиля нет", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			p := newPersonBuilder().
				withName().
				build()

			require.False(t, p.HasProfile())
			//Act
			initialUpdatedAt := p.UpdatedAt()
			time.Sleep(10 * time.Second)
			p.DetachProfile()
			_, ok := p.Profile()

			//Assert
			assert.False(t, p.HasProfile())
			assert.Equal(t, p.UpdatedAt(), initialUpdatedAt)
			assert.False(t, ok)
		})
	})

	t.Run("Rename обновляет имя и updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			p := newPersonBuilder().
				withName().
				build()

			//Act
			initialUpdatedAt := p.UpdatedAt()
			time.Sleep(1 * time.Second)
			newName, err := name.New(name.Params{
				Firstname:  "Петров",
				Lastname:   "Петр",
				Middlename: "Петрович",
			})
			require.NoError(t, err)

			p.Rename(newName)

			assert.Equal(t, newName, p.Name())
			assert.True(t, p.UpdatedAt().After(initialUpdatedAt))
		})
	})
}
