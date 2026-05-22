//go:build legacy
// +build legacy

package criteria_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Run("existing values returns true", func(t *testing.T) {
		//Assert
		assert.True(t, criteria.TypeRandom.IsValid())
		assert.True(t, criteria.TypeManual.IsValid())
	})

	t.Run("unexisting values returns false", func(t *testing.T) {
		//Assert
		assert.False(t, criteria.Type("").IsValid())
		assert.False(t, criteria.Type(" ").IsValid())
	})
}

func TestTitle(t *testing.T) {
	t.Run("for existing values returns non-empty string", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, criteria.TypeRandom.Title())
		assert.NotEmpty(t, criteria.TypeManual.Title())
	})

	t.Run("for non-existing values returns empty string", func(t *testing.T) {
		//Assert
		assert.Empty(t, criteria.Type("").Title())
		assert.Empty(t, criteria.Type(" ").Title())
	})
}
