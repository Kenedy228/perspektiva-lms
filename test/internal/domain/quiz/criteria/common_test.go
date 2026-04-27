package criteria_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/stretchr/testify/require"
)

func castManual(t *testing.T, c criteria.Criteria) criteria.Manual {
	cast, ok := c.(criteria.Manual)
	require.True(t, ok)

	return cast
}

func castRandom(t *testing.T, c criteria.Criteria) criteria.Random {
	cast, ok := c.(criteria.Random)
	require.True(t, ok)

	return cast
}
