package queries_test

import (
	"context"
	"errors"
	"testing"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/application/usecases/person/queries"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func shortViewsFixture() []personports.PersonShortView {
	return []personports.PersonShortView{{
		ID:               uuid.NewString(),
		FullName:         "Иванов Иван Иванович",
		OrganizationName: "ООО Ромашка",
	}}
}

func TestListByLastNameQuery_AdminPaginationAndNormalization(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListByLastNameQuery(&s)
	views := shortViewsFixture()

	s.On("ListByLastName", mock.Anything, "Иванов", common.DefaultLimit, 3).Return(views, nil)

	out, err := q.Execute(context.Background(), queries.ListByLastnameInput{
		ActorRole: role.NewAdmin(),
		LastName:  "  Иванов  ",
		Limit:     0,
		Offset:    3,
	})

	require.NoError(t, err)
	assert.Equal(t, views, out.Views)
	s.AssertExpectations(t)
}

func TestListByLastNameQuery_ForbiddenAndInvalidPagination(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListByLastNameQuery(&s)

	out, err := q.Execute(context.Background(), queries.ListByLastnameInput{
		ActorRole: role.NewStudent(),
		LastName:  "Иванов",
		Limit:     10,
	})
	assert.ErrorIs(t, err, common.ErrForbidden)
	assert.Nil(t, out)

	out, err = q.Execute(context.Background(), queries.ListByLastnameInput{
		ActorRole: role.NewAdmin(),
		LastName:  "Иванов",
		Limit:     -1,
	})
	assert.ErrorIs(t, err, common.ErrInvalidInput)
	assert.Nil(t, out)
}

func TestListBySnilsQuery_NormalizesSNILS(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListBySnilsQuery(&s)
	views := shortViewsFixture()

	s.On("ListBySnils", mock.Anything, "11223344595", 10, 0).Return(views, nil)

	out, err := q.Execute(context.Background(), queries.ListBySnilsInput{
		ActorRole: role.NewAdmin(),
		Snils:     "112-233-445 95",
		Limit:     10,
	})

	require.NoError(t, err)
	assert.Equal(t, views, out.Views)
}

func TestListByOrganizationIDQuery(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListByOrganizationIDQuery(&s)
	views := shortViewsFixture()
	orgID := uuid.New()

	s.On("ListByOrganizationID", mock.Anything, orgID, 10, 0).Return(views, nil)

	out, err := q.Execute(context.Background(), queries.ListByOrganizationIDInput{
		ActorRole:      role.NewAdmin(),
		OrganizationID: orgID.String(),
		Limit:          10,
	})

	require.NoError(t, err)
	assert.Equal(t, views, out.Views)
}

func TestGetDetailsByIDQuery_WrapsServiceError(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewGetDetailsByIdQuery(&s)
	id := uuid.New()

	s.On("GetDetailsByID", mock.Anything, id).Return(nil, errors.New("storage down"))

	out, err := q.Execute(context.Background(), queries.GetDetailsByIDInput{
		ActorRole: role.NewAdmin(),
		ID:        id.String(),
	})

	assert.Nil(t, out)
	assert.ErrorContains(t, err, "get person details")
}

func TestQueryConstructorsPanicOnNilDependencies(t *testing.T) {
	assert.Panics(t, func() { queries.NewGetDetailsByIdQuery(nil) })
	assert.Panics(t, func() { queries.NewListByLastNameQuery(nil) })
	assert.Panics(t, func() { queries.NewListByOrganizationIDQuery(nil) })
	assert.Panics(t, func() { queries.NewListBySnilsQuery(nil) })
}
