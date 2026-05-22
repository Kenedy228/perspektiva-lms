package queries_test

import (
	"context"
	"errors"
	"testing"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	"gitflic.ru/lms/backend/internal/application/usecases/account/queries"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func viewFixture(id uuid.UUID) accountports.AccountView {
	return accountports.AccountView{
		ID:       id.String(),
		Login:    "student2026",
		Role:     role.TypeStudent.String(),
		PersonID: uuid.NewString(),
		Status:   accountdomain.StatusActive.String(),
	}
}

func TestGetByIDQuery(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewGetByIDQuery(&s)
	id := uuid.New()
	view := viewFixture(id)

	s.On("GetByID", mock.Anything, id).Return(view, nil)

	out, err := q.Execute(context.Background(), queries.GetByIDInput{
		ActorRole: role.NewAdmin(),
		AccountID: id.String(),
	})

	require.NoError(t, err)
	assert.Equal(t, view, out.View)
}

func TestGetByPersonIDQuery_AdminOnly(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewGetByPersonIDQuery(&s)

	out, err := q.Execute(context.Background(), queries.GetByPersonIDInput{
		ActorRole: role.NewStudent(),
		PersonID:  uuid.NewString(),
	})

	assert.ErrorIs(t, err, common.ErrForbidden)
	assert.Nil(t, out)
}

func TestListQuery_NormalizesAndValidatesFilters(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListQuery(&s)
	views := []accountports.AccountView{viewFixture(uuid.New())}

	s.On("List", mock.Anything, accountports.ListFilter{
		Role:   role.TypeStudent,
		Status: accountdomain.StatusActive,
		Login:  "student 2026",
	}, common.DefaultLimit, 2).Return(views, nil)

	out, err := q.Execute(context.Background(), queries.ListInput{
		ActorRole: role.NewAdmin(),
		Role:      role.TypeStudent,
		Status:    accountdomain.StatusActive,
		Login:     "  student   2026  ",
		Offset:    2,
	})

	require.NoError(t, err)
	assert.Equal(t, views, out.Views)
}

func TestListQuery_InvalidPaginationAndFilters(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewListQuery(&s)

	out, err := q.Execute(context.Background(), queries.ListInput{
		ActorRole: role.NewAdmin(),
		Limit:     -1,
	})
	assert.ErrorIs(t, err, common.ErrInvalidInput)
	assert.Nil(t, out)

	out, err = q.Execute(context.Background(), queries.ListInput{
		ActorRole: role.NewAdmin(),
		Role:      role.Type("bad"),
	})
	assert.ErrorIs(t, err, common.ErrInvalidInput)
	assert.Nil(t, out)

	out, err = q.Execute(context.Background(), queries.ListInput{
		ActorRole: role.NewAdmin(),
		Status:    accountdomain.Status("bad"),
	})
	assert.ErrorIs(t, err, common.ErrInvalidInput)
	assert.Nil(t, out)
}

func TestQueryWrapsServiceError(t *testing.T) {
	s := mockQueryService{}
	q := queries.NewGetByIDQuery(&s)
	id := uuid.New()

	s.On("GetByID", mock.Anything, id).Return(nil, errors.New("storage down"))

	out, err := q.Execute(context.Background(), queries.GetByIDInput{
		ActorRole: role.NewAdmin(),
		AccountID: id.String(),
	})

	assert.Nil(t, out)
	assert.ErrorContains(t, err, "get account by id")
}

func TestConstructorsPanicOnNilDependencies(t *testing.T) {
	assert.Panics(t, func() { queries.NewGetByIDQuery(nil) })
	assert.Panics(t, func() { queries.NewGetByPersonIDQuery(nil) })
	assert.Panics(t, func() { queries.NewListQuery(nil) })
}
