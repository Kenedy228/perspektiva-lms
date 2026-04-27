package account_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount_Success(t *testing.T) {
	t.Run("создание с нормализацией логина", func(t *testing.T) {
		//Arrange
		acc := newAccountBuilder().
			withLogin("  User.Name  ").
			withPasswordHash("hashed-password").
			withRole(role.NewAdmin()).
			withPersonID(uuid.New()).
			build(t, nil)

		//Assert
		assert.NotEqual(t, uuid.Nil, acc.ID())
		assert.Equal(t, "user.name", acc.Login(), "логин должен быть нормализован")
		assert.Equal(t, "hashed-password", acc.PasswordHash())
		assert.False(t, acc.CreatedAt().IsZero())
		assert.False(t, acc.UpdatedAt().IsZero())
		assert.Equal(t, acc.CreatedAt(), acc.UpdatedAt(), "при создании даты должны совпадать")
	})

	t.Run("роль и personID сохраняются как есть", func(t *testing.T) {
		//Arrange
		personID := uuid.New()
		r := role.NewAdmin()

		acc := newAccountBuilder().
			withLogin("  User.Name  ").
			withPasswordHash("hashed-password").
			withRole(r).
			withPersonID(personID).
			build(t, nil)

		//Assert
		assert.Equal(t, r, acc.Role())
		assert.Equal(t, personID, acc.PersonID())
	})
}

func TestAccount_Mutations(t *testing.T) {
	t.Run("ChangeLogin нормализует и обновляет updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			acc := newAccountBuilder().
				withLogin("  User.Name  ").
				withPasswordHash("hashed-password").
				withRole(role.NewAdmin()).
				withPersonID(uuid.New()).
				build(t, nil)

			//Act
			initialUpdatedAt := acc.UpdatedAt()
			time.Sleep(10 * time.Second)
			err := acc.ChangeLogin("  New.Login  ")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, "new.login", acc.Login())
			assert.True(t, acc.UpdatedAt().After(initialUpdatedAt))
		})
	})

	t.Run("ChangePasswordHash обновляет хэш и updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			acc := newAccountBuilder().
				withLogin("  User.Name  ").
				withPasswordHash("hashed-password").
				withRole(role.NewAdmin()).
				withPersonID(uuid.New()).
				build(t, nil)

			//Act
			initialUpdatedAt := acc.UpdatedAt()
			time.Sleep(time.Second * 10)
			err := acc.ChangePasswordHash("new-hash")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, "new-hash", acc.PasswordHash())
			assert.True(t, acc.UpdatedAt().After(initialUpdatedAt))
		})
	})

	t.Run("ChangeRole обновляет роль и updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			acc := newAccountBuilder().
				withLogin("  User.Name  ").
				withPasswordHash("hashed-password").
				withRole(role.NewAdmin()).
				withPersonID(uuid.New()).
				build(t, nil)

			//Act
			initialUpdatedAt := acc.UpdatedAt()
			time.Sleep(time.Second * 10)
			newRole := role.NewCreator()
			acc.ChangeRole(newRole)

			//Assert
			assert.Equal(t, newRole, acc.Role())
			assert.True(t, acc.UpdatedAt().After(initialUpdatedAt))
		})
	})

	t.Run("ChangePersonID обновляет personID и updatedAt", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			acc := newAccountBuilder().
				withLogin("  User.Name  ").
				withPasswordHash("hashed-password").
				withRole(role.NewAdmin()).
				withPersonID(uuid.New()).
				build(t, nil)

			//Act
			initialUpdatedAt := acc.UpdatedAt()
			time.Sleep(time.Second * 10)
			newID := uuid.New()
			err := acc.ChangePersonID(newID)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, newID, acc.PersonID())
			assert.True(t, acc.UpdatedAt().After(initialUpdatedAt))
		})
	})
}
