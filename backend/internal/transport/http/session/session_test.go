package session

import (
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestManager_IssueAndVerify(t *testing.T) {
	manager := NewManager([]byte("secret"), time.Hour)
	accountID := uuid.New()
	personID := uuid.New()

	token, issued, err := manager.Issue(accountID, personID, role.NewAdmin())
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, accountID, issued.AccountID)

	verified, err := manager.Verify(token)
	require.NoError(t, err)
	require.Equal(t, accountID, verified.AccountID)
	require.Equal(t, personID, verified.PersonID)
	require.Equal(t, role.TypeAdmin, verified.Role.Kind())
}

func TestManager_VerifyRejectsTamperedToken(t *testing.T) {
	manager := NewManager([]byte("secret"), time.Hour)
	token, _, err := manager.Issue(uuid.New(), uuid.New(), role.NewStudent())
	require.NoError(t, err)

	_, err = manager.Verify(token + "tampered")
	require.ErrorIs(t, err, ErrInvalidToken)
}
