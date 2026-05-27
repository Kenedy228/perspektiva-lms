package short

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
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

func variantsSliceFixture(t *testing.T, count int) []variant.Variant {
	result := make([]variant.Variant, 0, count)

	for i := range count {
		val := fmt.Sprintf("val%d", i)
		v, err := variant.New(val)

		require.NoError(t, err)
		result = append(result, v)
	}

	return result
}
