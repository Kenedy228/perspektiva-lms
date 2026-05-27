package sequence

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"github.com/stretchr/testify/require"
)

func baseFixture(t *testing.T) *base.Base {
	t.Helper()

	titl, err := title.New("title")
	require.NoError(t, err)

	b, err := base.New(titl)

	require.NoError(t, err)

	return b
}

func optionsSliceFixture(t *testing.T, count int) []option.Option {
	t.Helper()

	result := make([]option.Option, 0, count)

	for i := range count {
		val := fmt.Sprintf("val%d", i)
		v, err := option.New(val)

		require.NoError(t, err)
		result = append(result, v)
	}

	return result
}
