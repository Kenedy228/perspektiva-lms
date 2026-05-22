package account

import (
	"context"

	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type QueryService interface {
	GetByID(ctx context.Context, id uuid.UUID) (AccountView, error)
	GetByPersonID(ctx context.Context, personID uuid.UUID) (AccountView, error)
	List(ctx context.Context, filter ListFilter, limit, offset int) ([]AccountView, error)
}

type ListFilter struct {
	Role   role.Type
	Status accountdomain.Status
	Login  string
}
