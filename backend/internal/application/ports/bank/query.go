package bank

import (
	"context"

	"github.com/google/uuid"
)

type Filter struct {
	TitleContains string
	QuestionID    uuid.UUID
	MinQuestions  int
	MaxQuestions  int
	Limit         int
	Offset        int
}

type QueryService interface {
	List(ctx context.Context, filter Filter) ([]ShortView, error)
	GetDetailsByID(ctx context.Context, id uuid.UUID) (DetailedView, error)
}
