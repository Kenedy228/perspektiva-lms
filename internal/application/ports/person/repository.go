package person

import (
	"context"

	"gitflic.ru/lms/internal/domain/person"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*person.Person, error)
	Save(ctx context.Context, p *person.Person) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
