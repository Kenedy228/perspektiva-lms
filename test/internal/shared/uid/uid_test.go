package uid_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("uid should generate id without error", func(t *testing.T) {
		//Arrange
		id, err := uid.New()
		
		//Assert
		assert.NoError(t, err)
		assert.NotEqual(t, id, uuid.Nil)
	})

	t.Run("uid should parse from string to uuid", func(t *testing.T) {
		//Arrange
		id, _ := uid.New()
		
		//Act
		_, err := uuid.Parse(id.String())
		
		//Assert
		assert.NoError(t, err)
	})

	t.Run("different New calls generate different ids", func(t *testing.T) {
		//Arrange
		first, _ := uid.New()
		second, _ := uid.New()

		//Assert
		assert.NotEqual(t, first, second)
	})
}

