package question

import (
	"context"

	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (questdomain.Question, error)
	Save(ctx context.Context, q questdomain.Question) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
