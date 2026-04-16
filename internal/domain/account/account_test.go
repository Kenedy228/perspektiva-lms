package account

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var mockRole role.Role

	tests := []struct {
		name         string
		login        string
		passwordHash string
		personID     uuid.UUID
		err          error
	}{
		{
			name:         "success valid account",
			login:        "user@example.com",
			passwordHash: "$2a$12$somehashstring",
			personID:     uuid.New(),
			err:          nil,
		},
		{
			name:         "error empty login",
			login:        "",
			passwordHash: "hash",
			personID:     uuid.New(),
			err:          ErrEmptyLogin,
		},
		{
			name:         "error whitespaces login",
			login:        "   ",
			passwordHash: "hash",
			personID:     uuid.New(),
			err:          ErrEmptyLogin,
		},
		{
			name:         "error empty password hash",
			login:        "user",
			passwordHash: "",
			personID:     uuid.New(),
			err:          ErrEmptyPasswordHash,
		},
		{
			name:         "error nil person id",
			login:        "user",
			passwordHash: "hash",
			personID:     uuid.Nil,
			err:          ErrNilPersonID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := Params{
				Login:        tt.login,
				PasswordHash: tt.passwordHash,
				Role:         mockRole,
				PersonID:     tt.personID,
			}

			acc, err := New(params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, acc)
				assert.NotEqual(t, uuid.Nil, acc.ID())
				assert.Equal(t, tt.login, acc.Login())
				assert.Equal(t, tt.personID, acc.PersonID())
				assert.False(t, acc.IsBlocked(), "new account should not be blocked")
				assert.False(t, acc.CreatedAt().IsZero())
				assert.False(t, acc.UpdatedAt().IsZero())
				assert.Equal(t, acc.CreatedAt(), acc.UpdatedAt())
			}
		})
	}
}

func TestAccount_ChangeLogin(t *testing.T) {
	acc := baseAccount()
	oldUpdatedAt := acc.UpdatedAt()

	time.Sleep(time.Millisecond * 5) // Небольшая пауза, чтобы время точно изменилось

	t.Run("success", func(t *testing.T) {
		err := acc.ChangeLogin("new_login")

		assert.NoError(t, err)
		assert.Equal(t, "new_login", acc.Login())
		assert.True(t, acc.UpdatedAt().After(oldUpdatedAt))
	})

	t.Run("error empty login keeps old state", func(t *testing.T) {
		currentLogin := acc.Login()
		currentUpdateAt := acc.UpdatedAt()

		err := acc.ChangeLogin("   ")

		assert.ErrorIs(t, err, ErrEmptyLogin)
		assert.Equal(t, currentLogin, acc.Login())
		assert.Equal(t, currentUpdateAt, acc.UpdatedAt())
	})
}

func TestAccount_ChangePassword(t *testing.T) {
	acc := baseAccount()

	t.Run("success", func(t *testing.T) {
		err := acc.ChangePassword("new_hash")
		assert.NoError(t, err)
		// Убедиться напрямую мы не можем (нет геттера для хэша, и это правильно для безопасности),
		// но мы знаем, что ошибка не вернулась.
	})

	t.Run("error empty hash", func(t *testing.T) {
		err := acc.ChangePassword("")
		assert.ErrorIs(t, err, ErrEmptyPasswordHash)
	})
}

func TestAccount_ChangeRole(t *testing.T) {
	acc := baseAccount()
	var newRole role.Role // В реальном тесте подставьте валидную роль

	acc.ChangeRole(newRole)
	assert.Equal(t, newRole, acc.Role())
}

func TestAccount_ChangePersonID(t *testing.T) {
	acc := baseAccount()
	oldUpdatedAt := acc.UpdatedAt()
	time.Sleep(time.Millisecond * 5)

	t.Run("success", func(t *testing.T) {
		newID := uuid.New()
		err := acc.ChangePersonID(newID)

		assert.NoError(t, err)
		assert.Equal(t, newID, acc.PersonID())
		assert.True(t, acc.UpdatedAt().After(oldUpdatedAt))
	})

	t.Run("error nil id", func(t *testing.T) {
		currentID := acc.PersonID()

		err := acc.ChangePersonID(uuid.Nil)

		assert.ErrorIs(t, err, ErrNilPersonID)
		assert.Equal(t, currentID, acc.PersonID())
	})
}

func TestAccount_BlockUnblock(t *testing.T) {
	acc := baseAccount()

	// Начальное состояние: не заблокирован
	assert.False(t, acc.IsBlocked())
	initialUpdatedAt := acc.UpdatedAt()

	time.Sleep(time.Millisecond * 5)

	t.Run("first block changes state and time", func(t *testing.T) {
		acc.Block()

		assert.True(t, acc.IsBlocked())
		assert.True(t, acc.UpdatedAt().After(initialUpdatedAt))
	})

	t.Run("second block is idempotent", func(t *testing.T) {
		blockedUpdatedAt := acc.UpdatedAt()
		time.Sleep(time.Millisecond * 5)

		acc.Block() // Повторная блокировка

		assert.True(t, acc.IsBlocked()) // Состояние осталось прежним
		assert.Equal(t, blockedUpdatedAt, acc.UpdatedAt(), "updatedAt should not change on redundant block")
	})

	t.Run("first unblock changes state and time", func(t *testing.T) {
		blockedUpdatedAt := acc.UpdatedAt()
		time.Sleep(time.Millisecond * 5)

		acc.Unblock()

		assert.False(t, acc.IsBlocked())
		assert.True(t, acc.UpdatedAt().After(blockedUpdatedAt))
	})

	t.Run("second unblock is idempotent", func(t *testing.T) {
		unblockedUpdatedAt := acc.UpdatedAt()
		time.Sleep(time.Millisecond * 5)

		acc.Unblock() // Повторная разблокировка

		assert.False(t, acc.IsBlocked()) // Состояние осталось прежним
		assert.Equal(t, unblockedUpdatedAt, acc.UpdatedAt(), "updatedAt should not change on redundant unblock")
	})
}

func TestAccount_Equal(t *testing.T) {
	id := uuid.New()

	acc1 := &Account{id: id}
	acc2 := &Account{id: id}
	acc3 := &Account{id: uuid.New()}

	assert.True(t, acc1.Equal(acc2), "accounts with same id should be equal")
	assert.False(t, acc1.Equal(acc3), "accounts with different ids should not be equal")
	assert.False(t, acc1.Equal(nil), "account should not be equal to nil")
}

// === Вспомогательные функции ===

func baseAccount() *Account {
	var mockRole role.Role

	acc, _ := New(Params{
		Login:        "admin",
		PasswordHash: "hash",
		Role:         mockRole,
		PersonID:     uuid.New(),
	})

	return acc
}
