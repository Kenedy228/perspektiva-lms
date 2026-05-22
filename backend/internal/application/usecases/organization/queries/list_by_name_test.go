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

func TestListByNameQuery(t *testing.T) {
	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := queries.ListByNameInput{
			ActorRole: studentRole(),
			Name:      "ООО 'Ромашка'",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByNameQuery(&s)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		s.AssertNotCalled(t, "ListByName", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("некорректная пагинация", func(t *testing.T) {
		in := queries.ListByNameInput{
			ActorRole: adminRole(),
			Name:      "ООО 'Ромашка'",
			Limit:     common.MaxLimit + 1,
			Offset:    0,
		}

		s := mockQueryService{}
		uc := queries.NewListByNameQuery(&s)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Nil(t, out)
		s.AssertNotCalled(t, "ListByName", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ошибка запроса на выборку", func(t *testing.T) {
		// Arrange
		in := queries.ListByNameInput{
			ActorRole: adminRole(),
			Name:      "ООО 'Ромашка'",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByNameQuery(&s)
		s.On("ListByName", mock.Anything, "ООО 'Ромашка'", 10, 20).Return(nil, errors.New(""))

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
		in := queries.ListByNameInput{
			ActorRole: adminRole(),
			Name:      "ООО 'Ромашка'",
			Limit:     10,
			Offset:    20,
		}

		s := mockQueryService{}
		uc := queries.NewListByNameQuery(&s)
		s.On("ListByName", mock.Anything, "ООО 'Ромашка'", 10, 20).Return(orgs, nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, orgs, out.Views)
		s.AssertExpectations(t)
	})

	t.Run("нормализует строку поиска и использует лимит по умолчанию", func(t *testing.T) {
		orgs := organizationShortSliceViewFixture()
		in := queries.ListByNameInput{
			ActorRole: adminRole(),
			Name:      "  ООО   'Ромашка'  ",
			Limit:     0,
			Offset:    5,
		}

		s := mockQueryService{}
		uc := queries.NewListByNameQuery(&s)
		s.On("ListByName", mock.Anything, "ООО 'Ромашка'", common.DefaultLimit, 5).Return(orgs, nil)

		out, err := uc.Execute(context.Background(), in)

		assert.NoError(t, err)
		assert.Equal(t, orgs, out.Views)
		s.AssertExpectations(t)
	})
}
