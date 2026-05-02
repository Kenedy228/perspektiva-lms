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

func TestDetachProfileUseCase(t *testing.T) {
	t.Run("ошибка при поиске профиля", func(t *testing.T) {
		// Arrange
		in := commands.DetachProfileInput{
			PersonID: uuid.NewString(),
		}

		r := mockRepository{}
		uc := commands.NewDetachProfileUseCase(&r)
		r.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("при возврате записи без профиля Save() не вызывается", func(t *testing.T) {
		// Arrange
		pFixture := personFixture()
		in := commands.DetachProfileInput{
			PersonID: pFixture.ID().String(),
		}

		r := mockRepository{}
		uc := commands.NewDetachProfileUseCase(&r)
		r.On("FindByID", mock.Anything, mock.Anything).Return(pFixture, nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, pFixture.ID().String(), out.PersonID)
		r.AssertExpectations(t)
	})

	t.Run("при возврате записи с профилем Save() вызывается", func(t *testing.T) {
		// Arrange
		pFixture := personWithProfileFixture()
		in := commands.DetachProfileInput{
			PersonID: pFixture.ID().String(),
		}

		r := mockRepository{}
		uc := commands.NewDetachProfileUseCase(&r)
		r.On("FindByID", mock.Anything, mock.Anything).Return(pFixture, nil)
		r.On("Save", mock.Anything, pFixture).Return(nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, pFixture.ID().String(), out.PersonID)
		r.AssertExpectations(t)
	})
}
