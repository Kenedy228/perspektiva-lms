//go:build legacy
// +build legacy

package short_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/short"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Success(t *testing.T) {
	//Arrange
	q, err := short.New(makeTitle(), makeVariants(10))

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, 10, len(q.Variants()))
	assert.Equal(t, question.TypeShort.DefaultInstruction(), q.Instruction())
	assert.Equal(t, question.TypeShort, q.Type())
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name          string
		variantsCount int
		wantErr       error
	}{
		{
			name:          "количество вариантов меньше необходимого",
			variantsCount: 0,
			wantErr:       short.ErrInvalid,
		},
		{
			name:          "количество вариантов больше необходимого",
			variantsCount: 1e2,
			wantErr:       short.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Assert
			_, err := short.New(makeTitle(), makeVariants(tt.variantsCount))

			//Arrange
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestVariants(t *testing.T) {
	//Arrange
	q, err := short.New(makeTitle(), makeVariants(10))
	require.NoError(t, err)
	variants := q.Variants()

	//Act
	variants = makeVariants(1)

	//Assert
	assert.NotEqual(t, variants, q.Variants())
}

func TestChangeVariants_Success(t *testing.T) {
	//Arrange
	q, err := short.New(makeTitle(), makeVariants(10))
	require.NoError(t, err)
	newVariants := makeVariants(10)

	//Act
	err = q.ChangeVariants(newVariants)

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, newVariants, q.Variants())
}

func TestChangeVariants_Fail(t *testing.T) {
	tc := []struct {
		name          string
		variantsCount int
		wantErr       error
	}{
		{
			name:          "количество вариантов меньше необходимого",
			variantsCount: 0,
			wantErr:       short.ErrInvalid,
		},
		{
			name:          "количество вариантов больше необходимого",
			variantsCount: 1e2,
			wantErr:       short.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			q, err := short.New(makeTitle(), makeVariants(10))
			require.NoError(t, err)

			//Act
			err = q.ChangeVariants(makeVariants(tt.variantsCount))

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	q, err := short.New(makeTitle(), makeVariants(10))
	require.NoError(t, err)
	clone, ok := q.Clone().(*short.Question)
	require.True(t, ok)

	//Assert
	assert.Equal(t, q.Variants(), clone.Variants())
}
