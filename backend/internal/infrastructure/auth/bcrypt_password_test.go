package auth_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptPasswordComparer_Compare(t *testing.T) {
	comparer, err := auth.NewBcryptPasswordComparerWithCost(4)
	require.NoError(t, err)

	hash, err := comparer.Hash("password")
	require.NoError(t, err)

	tests := []struct {
		name  string
		hash  passhash.Hash
		plain string
		want  bool
	}{
		{
			name:  "returns true when password matches hash",
			hash:  hash,
			plain: "password",
			want:  true,
		},
		{
			name:  "returns false when password does not match hash",
			hash:  hash,
			plain: "wrong-password",
			want:  false,
		},
		{
			name:  "returns false for empty hash",
			hash:  passhash.Hash{},
			plain: "password",
			want:  false,
		},
		{
			name:  "returns false for non-bcrypt hash",
			hash:  mustHash(t, "argon2id:v=19:m=65536,t=3,p=2:salt:hash"),
			plain: "password",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, comparer.Compare(tt.hash, tt.plain))
		})
	}
}

func TestBcryptPasswordComparer_Hash(t *testing.T) {
	comparer, err := auth.NewBcryptPasswordComparerWithCost(4)
	require.NoError(t, err)

	hash, err := comparer.Hash("new-password")

	require.NoError(t, err)
	assert.NotEmpty(t, hash.Value())
	assert.True(t, comparer.Compare(hash, "new-password"))
	assert.False(t, comparer.Compare(hash, "wrong-password"))
}

func TestNewBcryptPasswordComparerWithCost(t *testing.T) {
	_, err := auth.NewBcryptPasswordComparerWithCost(3)

	require.Error(t, err)
}

func mustHash(t *testing.T, value string) passhash.Hash {
	t.Helper()

	hash, err := passhash.New(value)
	require.NoError(t, err)

	return hash
}
