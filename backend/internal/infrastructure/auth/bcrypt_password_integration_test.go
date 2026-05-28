package auth_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/infrastructure/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptPasswordComparer_AccountIntegration(t *testing.T) {
	comparer, err := auth.NewBcryptPasswordComparerWithCost(10)
	require.NoError(t, err)

	hash, err := comparer.Hash("student-password")
	require.NoError(t, err)

	l, err := login.New("student26")
	require.NoError(t, err)

	personID := uuid.New()

	acc, err := account.New(l, hash, role.NewStudent(), personID)
	require.NoError(t, err)

	assert.True(t, comparer.Compare(acc.PasswordHash(), "student-password"))
	assert.False(t, comparer.Compare(acc.PasswordHash(), "wrong-password"))
}
