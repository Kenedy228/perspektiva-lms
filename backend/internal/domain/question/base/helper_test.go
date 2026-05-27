package base

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var idFixture = uuid.New()

func titleFixture(t *testing.T) title.Title {
	titl, err := title.New("value")
	require.NoError(t, err)

	return titl
}
