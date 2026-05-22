package queries_test

import (
	"context"
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListByINNQuery(t *testing.T) {
	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := queries.ListByINNInput{
			ActorRole: studentRole(),
			INN:       "0000000000",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByINNQuery(&s)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		s.AssertNotCalled(t, "ListByINN", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("пустая строка поиска", func(t *testing.T) {
		in := queries.ListByINNInput{
			ActorRole: adminRole(),
			INN:       "   ",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByINNQuery(&s)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Nil(t, out)
		s.AssertNotCalled(t, "ListByINN", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ошибка запроса на выборку", func(t *testing.T) {
		// Arrange
		in := queries.ListByINNInput{
			ActorRole: adminRole(),
			INN:       "0000000000",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByINNQuery(&s)
		s.On("ListByINN", mock.Anything, "0000000000", 10, 20).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		s.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		orgs := organizationShortSliceViewFixture()
		in := queries.ListByINNInput{
			ActorRole: adminRole(),
			INN:       "0000000000",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByINNQuery(&s)
		s.On("ListByINN", mock.Anything, "0000000000", 10, 20).Return(orgs, nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, orgs, out.Views)
		s.AssertExpectations(t)
	})

	t.Run("нормализует ИНН и использует лимит по умолчанию", func(t *testing.T) {
		orgs := organizationShortSliceViewFixture()
		in := queries.ListByINNInput{
			ActorRole: adminRole(),
			INN:       "  000 000 0000 ",
			Limit:     0,
			Offset:    3,
		}

		s := mockQueryService{}
		uc := queries.NewListByINNQuery(&s)
		s.On("ListByINN", mock.Anything, "0000000000", common.DefaultLimit, 3).Return(orgs, nil)

		out, err := uc.Execute(context.Background(), in)

		assert.NoError(t, err)
		assert.Equal(t, orgs, out.Views)
		s.AssertExpectations(t)
	})
}
