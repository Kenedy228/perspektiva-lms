package queries_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetDetailsByIDQuery(t *testing.T) {
	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := queries.GetDetailsByIDInput{
			ActorRole: studentRole(),
			ID:        uuid.NewString(),
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIDQuery(&s)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		s.AssertNotCalled(t, "GetDetailsByID", mock.Anything, mock.Anything)
	})

	t.Run("некорректный идентификатор", func(t *testing.T) {
		// Arrange
		in := queries.GetDetailsByIDInput{
			ActorRole: adminRole(),
			ID:        "xxx",
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIDQuery(&s)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		s.AssertExpectations(t)
	})

	t.Run("ошибка запроса на выборку", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		in := queries.GetDetailsByIDInput{
			ActorRole: adminRole(),
			ID:        id.String(),
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIDQuery(&s)
		s.On("GetDetailsByID", mock.Anything, id).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		s.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		org := organizationDetailedViewFixture()
		id, err := uuid.Parse(org.ID)
		require.NoError(t, err)

		in := queries.GetDetailsByIDInput{
			ActorRole: adminRole(),
			ID:        id.String(),
		}

		s := mockQueryService{}
		uc := queries.NewGetDetailsByIDQuery(&s)
		s.On("GetDetailsByID", mock.Anything, id).Return(org, nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, org, out.View)
		s.AssertExpectations(t)
	})
}
