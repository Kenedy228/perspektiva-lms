package quiz

import (
	"context"

	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*quizdomain.Quiz, error)
	Save(ctx context.Context, q *quizdomain.Quiz) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type QuestionBankInspector interface {
	CountQuestionsInBank(ctx context.Context, bankID uuid.UUID) (int, error)
	QuestionsBelongToBank(ctx context.Context, bankID uuid.UUID, questionIDs []uuid.UUID) (bool, error)
}
