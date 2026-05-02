package queries_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/internal/application/usecases/person/queries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDetailsByID(t *testing.T) {
	t.Run("ошибка при парсинге идентификатора", func(t *testing.T) {
		// Arrange
		in := queries.GetDetailsByIDInput{
			ID: "",
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIdQuery(&s)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		s.AssertExpectations(t)
	})

	t.Run("ошибка при поиске значения", func(t *testing.T) {
		// Arrange
		in := queries.GetDetailsByIDInput{
			ID: uuid.NewString(),
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIdQuery(&s)
		s.On("GetDetailsByID", mock.Anything, mock.Anything).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		s.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		viewFixture := personDetailedFixture()
		in := queries.GetDetailsByIDInput{
			ID: viewFixture.ID,
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIdQuery(&s)
		s.On("GetDetailsByID", mock.Anything, mock.Anything).Return(viewFixture, nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, viewFixture, out.View)
		s.AssertExpectations(t)
	})
}
