package account_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/account"
	"gitflic.ru/lms/internal/domain/account/login"
	"gitflic.ru/lms/internal/domain/account/passhash"
	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("создание с personID", func(t *testing.T) {
		//Arrange
		personID := uuid.New()
		acc, err := newAccountBuilder().
			withLogin().
			withPasswordHash().
			withRole().
			withPersonID(personID).
			build()

		//Assert
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, acc.ID())
		assert.Equal(t, acc.Login(), loginFixture())
		assert.Equal(t, acc.PasswordHash(), hashFixture())
		assert.Equal(t, acc.Role(), roleFixture())
		assert.Equal(t, acc.PersonID(), personID)
	})

	t.Run("возвращает ошибку, если не установлен personID", func(t *testing.T) {
		//Arrange
		_, err := newAccountBuilder().
			withLogin().
			withPasswordHash().
			withRole().
			withPersonID(uuid.Nil).
			build()

		//Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, account.ErrInvalid)
	})
}

func TestChangeLogin(t *testing.T) {
	//Arrange
	acc, err := newAccountBuilder().
		withLogin().
		withPasswordHash().
		withRole().
		withPersonID(uuid.New()).
		build()
	require.NoError(t, err)

	//Act
	newLogin, _ := login.New("ivan.ivanov-99")
	acc.ChangeLogin(newLogin)

	//Assert
	assert.Equal(t, acc.Login(), newLogin)
}

func TestChangePasswordHash(t *testing.T) {
	//Arrange
	acc, err := newAccountBuilder().
		withLogin().
		withPasswordHash().
		withRole().
		withPersonID(uuid.New()).
		build()
	require.NoError(t, err)

	//Act
	newHash, _ := passhash.New("$2b$12$D24p4h.P6P4.82.jR.X1U.1Q6Qx4/G0iB2.JzH8H4w2rP/T5k0eZ2")
	acc.ChangePasswordHash(newHash)

	//Assert
	assert.Equal(t, acc.PasswordHash(), newHash)
}

func TestChangeRole(t *testing.T) {
	//Arrange
	acc, err := newAccountBuilder().
		withLogin().
		withPasswordHash().
		withRole().
		withPersonID(uuid.New()).
		build()
	require.NoError(t, err)

	//Act
	newRole := role.NewCreator()
	acc.ChangeRole(newRole)

	//Assert
	assert.Equal(t, acc.Role(), newRole)
}
