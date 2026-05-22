package account

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func loginFixture(t *testing.T) login.Login {
	l, err := login.New("login")
	require.NoError(t, err)
	return l
}

func passhashFixture(t *testing.T) passhash.Hash {
	h, err := passhash.New("hash")
	require.NoError(t, err)
	return h
}

var idFixture = uuid.New()
