package person

import (
	"context"

	"github.com/google/uuid"
)

type QueryService interface {
	ListByOrganizationID(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]PersonShortView, error)
	ListByLastName(ctx context.Context, lastName string, limit, offset int) ([]PersonShortView, error)
	ListBySnils(ctx context.Context, snils string, limit, offset int) ([]PersonShortView, error)
	GetDetailsByID(ctx context.Context, id uuid.UUID) (PersonDetailedView, error)
}
