//go:build legacy
// +build legacy

package criteria_test

import (
	"testing"

	criteria2 "gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/stretchr/testify/require"
)

func castManual(t *testing.T, c criteria2.Criteria) criteria2.Manual {
	cast, ok := c.(criteria2.Manual)
	require.True(t, ok)

	return cast
}

func castRandom(t *testing.T, c criteria2.Criteria) criteria2.Random {
	cast, ok := c.(criteria2.Random)
	require.True(t, ok)

	return cast
}
