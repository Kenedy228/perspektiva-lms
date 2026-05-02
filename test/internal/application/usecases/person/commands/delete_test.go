package commands_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/internal/application/usecases/person/commands"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUseCase(t *testing.T) {
	t.Run("ошибка при парсинге идентификатора", func(t *testing.T) {
		// Arrange
		in := commands.DeleteByIDInput{
			PersonID: "",
		}

		r := mockRepository{}
		uc := commands.NewDeleteByIDUseCase(&r)

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("ошибка при удалении записи", func(t *testing.T) {
		// Arrange
		pID := uuid.New()
		in := commands.DeleteByIDInput{
			PersonID: pID.String(),
		}

		r := mockRepository{}
		uc := commands.NewDeleteByIDUseCase(&r)
		r.On("DeleteByID", mock.Anything, pID).Return(errors.New(""))

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		pID := uuid.New()
		in := commands.DeleteByIDInput{
			PersonID: pID.String(),
		}

		r := mockRepository{}
		uc := commands.NewDeleteByIDUseCase(&r)
		r.On("DeleteByID", mock.Anything, pID).Return(nil)

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		r.AssertExpectations(t)
	})
}
