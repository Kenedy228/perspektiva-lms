//go:build legacy
// +build legacy

package source_test

import (
	"testing"

	source2 "gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSource(t *testing.T) {
	validBankID := uuid.New()
	validCriteria := criteriaFixture{}

	t.Run("успешное создание", func(t *testing.T) {
		// Arrange & Act
		s, err := source2.NewSource(validBankID, validCriteria)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, validBankID, s.BankID())
		assert.NotNil(t, s.Criteria())
	})

	t.Run("ошибка при пустом UUID банка (uuid.Nil)", func(t *testing.T) {
		// Arrange & Act
		_, err := source2.NewSource(uuid.Nil, validCriteria)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, source2.ErrInvalid)
		assert.Contains(t, err.Error(), "несуществующий банк вопросов")
	})

	t.Run("ошибка при nil criteria", func(t *testing.T) {
		// Arrange
		// Важно: мы передаем явный nil интерфейс
		var nilCriteria criteria.Criteria = nil

		// Act
		_, err := source2.NewSource(validBankID, nilCriteria)

		// Assert
		require.Error(t, err)
		assert.ErrorIs(t, err, source2.ErrInvalid)
		assert.Contains(t, err.Error(), "не выбран критерий выборки")
	})
}
