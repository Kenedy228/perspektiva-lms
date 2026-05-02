package commands_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/internal/application/usecases/person/commands"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePersonUseCase(t *testing.T) {
	t.Run("ошибка при создании объекта Name", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			FirstName:  "",
			LastName:   "",
			MiddleName: "",
		}

		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)
		require.NotNil(t, uc)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})

	t.Run("ошибка при сохранении", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			FirstName:  "Иван",
			LastName:   "Иванов",
			MiddleName: "Петрович",
		}

		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)
		require.NotNil(t, uc)
		r.On("Save", mock.Anything, mock.AnythingOfType("*person.Person")).Return(errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			FirstName:  "Иван",
			LastName:   "Иванов",
			MiddleName: "Петрович",
		}

		r := mockRepository{}
		uc := commands.NewCreateUseCase(&r)
		require.NotNil(t, uc)

		r.On("Save", mock.Anything, mock.AnythingOfType("*person.Person")).Return(nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, out.ID)
		r.AssertExpectations(t)
	})
}
