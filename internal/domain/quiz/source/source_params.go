package source

import (
	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
)

type Params struct {
	BankID uuid.UUID
	Criteria criteria.Criteria
}
