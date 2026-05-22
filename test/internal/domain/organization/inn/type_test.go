//go:build legacy
// +build legacy

package inn_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"github.com/stretchr/testify/assert"
)

func TestTypeTitle(t *testing.T) {
	t.Run("для определенных констант возвращает непустую строку", func(t *testing.T) {
		// Assert
		assert.NotEmpty(t, inn.TypeIP.Title())
		assert.NotEmpty(t, inn.TypeOrganization.Title())
		assert.NotEmpty(t, inn.TypePhysical.Title())
	})

	t.Run("для неопределенных типов возвращает пустую строку", func(t *testing.T) {
		// Assert
		assert.Empty(t, inn.Type("undefined").Title())
	})
}

func TestTypeCodeLength(t *testing.T) {
	t.Run("для определенных констант возвращает ненулевое значение", func(t *testing.T) {
		// Assert
		assert.NotZero(t, inn.TypeIP.CodeLength())
		assert.NotZero(t, inn.TypeOrganization.CodeLength())
		assert.NotZero(t, inn.TypePhysical.CodeLength())
	})

	t.Run("для неопределенных типов возвращает нулевое значение", func(t *testing.T) {
		// Assert
		assert.Zero(t, inn.Type("undefined").CodeLength())
	})
}

func TestTypeCoefficients(t *testing.T) {
	t.Run("для определенных констант возвращает непустой слайс", func(t *testing.T) {
		// Assert
		assert.NotEmpty(t, inn.TypeIP.Coefficients())
		assert.NotEmpty(t, inn.TypeOrganization.Coefficients())
		assert.NotEmpty(t, inn.TypePhysical.Coefficients())
	})

	t.Run("для неопределенных констант возвращает пустой слайс", func(t *testing.T) {
		// Assert
		assert.Empty(t, inn.Type("undefined").Coefficients())
	})
}
