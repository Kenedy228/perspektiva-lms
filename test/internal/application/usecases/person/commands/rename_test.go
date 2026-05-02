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

func TestRenamePersonUseCase(t *testing.T) {
	t.Run("ошибка при создании имени", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			PersonID:   uuid.NewString(),
			FirstName:  "",
			LastName:   "",
			MiddleName: "",
		}

		r := mockRepository{}
		uc := commands.NewRenameUseCase(&r)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Nil(t, out)
		assert.Error(t, err)
		r.AssertNotCalled(t, "FindByID")
		r.AssertNotCalled(t, "Save")
	})

	t.Run("ошибка при парсинге uuid", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			PersonID:   "uuid",
			FirstName:  "",
			LastName:   "",
			MiddleName: "",
		}

		r := mockRepository{}
		uc := commands.NewRenameUseCase(&r)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Nil(t, out)
		assert.Error(t, err)
		r.AssertNotCalled(t, "FindByID")
		r.AssertNotCalled(t, "Save")
	})

	t.Run("ошибка при поиске записи по ID", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			PersonID:   uuid.NewString(),
			FirstName:  "Иванов",
			LastName:   "Иван",
			MiddleName: "Иванович",
		}

		r := mockRepository{}
		uc := commands.NewRenameUseCase(&r)
		r.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Nil(t, out)
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("ошибка сохраенении изменений", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			PersonID:   uuid.NewString(),
			FirstName:  "Иванов",
			LastName:   "Иван",
			MiddleName: "Иванович",
		}

		r := mockRepository{}
		uc := commands.NewRenameUseCase(&r)
		r.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(personFixture(), nil)
		r.On("Save", mock.Anything, mock.AnythingOfType("*person.Person")).Return(errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Nil(t, out)
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		pFixture := personFixture()
		in := commands.RenameInput{
			PersonID:   pFixture.ID().String(),
			FirstName:  "Иванов",
			LastName:   "Иван",
			MiddleName: "Иванович",
		}

		r := mockRepository{}
		uc := commands.NewRenameUseCase(&r)
		r.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(pFixture, nil)
		r.On("Save", mock.Anything, mock.AnythingOfType("*person.Person")).Return(nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, pFixture.ID().String(), out.ID)
		r.AssertExpectations(t)
	})
}
