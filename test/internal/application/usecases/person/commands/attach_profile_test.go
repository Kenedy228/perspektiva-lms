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

func TestAttachProfileUseCase(t *testing.T) {
	t.Run("ошибка при создании профиля", func(t *testing.T) {
		// Arrange
		in := commands.AttachProfileInput{
			DateOfBirth:    dateOfBirthFixture().Date(),
			PersonID:       uuid.NewString(),
			Snils:          snilsFixture().Value(),
			JobTitle:       "",
			Education:      educationFixture().Value(),
			OrganizationID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewAttachProfileUseCase(&r)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка при парсинге PersonID", func(t *testing.T) {
		// Arrange
		in := commands.AttachProfileInput{
			DateOfBirth:    dateOfBirthFixture().Date(),
			PersonID:       "",
			Snils:          snilsFixture().Value(),
			JobTitle:       jobTitleFixture().Title(),
			Education:      educationFixture().Value(),
			OrganizationID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewAttachProfileUseCase(&r)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка при поиске записи", func(t *testing.T) {
		// Arrange
		in := commands.AttachProfileInput{
			DateOfBirth:    dateOfBirthFixture().Date(),
			PersonID:       uuid.NewString(),
			Snils:          snilsFixture().Value(),
			JobTitle:       jobTitleFixture().Title(),
			Education:      educationFixture().Value(),
			OrganizationID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewAttachProfileUseCase(&r)
		r.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка при сохранении записи", func(t *testing.T) {
		// Arrange
		pFixture := personFixture()

		in := commands.AttachProfileInput{
			DateOfBirth:    dateOfBirthFixture().Date(),
			PersonID:       pFixture.ID().String(),
			Snils:          snilsFixture().Value(),
			JobTitle:       jobTitleFixture().Title(),
			Education:      educationFixture().Value(),
			OrganizationID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewAttachProfileUseCase(&r)
		r.On("FindByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(pFixture, nil)
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
		pFixture := personFixture()

		in := commands.AttachProfileInput{
			DateOfBirth:    dateOfBirthFixture().Date(),
			PersonID:       pFixture.ID().String(),
			Snils:          snilsFixture().Value(),
			JobTitle:       jobTitleFixture().Title(),
			Education:      educationFixture().Value(),
			OrganizationID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewAttachProfileUseCase(&r)
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
