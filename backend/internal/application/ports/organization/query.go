package organization

import (
	"context"

	"github.com/google/uuid"
)

type QueryService interface {
	ListByName(ctx context.Context, name string, limit, offset int) ([]OrganizationShortView, error)
	ListByINN(ctx context.Context, inn string, limit, offset int) ([]OrganizationShortView, error)
	GetDetailsByID(ctx context.Context, id uuid.UUID) (OrganizationDetailedView, error)
}
