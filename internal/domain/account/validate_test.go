package account

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateLogin(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		//Act
		err := validateLogin("valid_login")

		//Assert
		assert.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		//Act
		err := validateLogin("   ")

		//Assert
		assert.ErrorIs(t, err, ErrInvalid)
	})

	t.Run("too long", func(t *testing.T) {
		//Arrange
		long := strings.Repeat("A", loginCharsLimit+500)

		//Act
		err := validateLogin(string(long))

		//Assert
		assert.ErrorIs(t, err, ErrInvalid)
	})
}

func TestValidatePasswordHash(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		//Act
		err := validatePasswordHash("hash")

		//Assert
		assert.NoError(t, err)
	})

	t.Run("empty", func(t *testing.T) {
		//Act
		err := validatePasswordHash("   ")

		//Assert
		assert.ErrorIs(t, err, ErrInvalid)
	})
}

func TestValidatePersonID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		//Act
		err := validatePersonID(uuid.New())

		//Assert
		assert.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		//Act
		err := validatePersonID(uuid.Nil)

		//Assert
		assert.ErrorIs(t, err, ErrInvalid)
	})
}

func TestNormalize(t *testing.T) {
	t.Run("trims spaces, removes inner spaces, lowercases", func(t *testing.T) {
		//Act
		got := normalize("  Us er Name  ")

		//Assert
		assert.Equal(t, "username", got)
	})
}
