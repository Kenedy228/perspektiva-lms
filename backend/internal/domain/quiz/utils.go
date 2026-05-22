package quiz

import (
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"github.com/google/uuid"
)

func getBankIDs(sources []source.Source) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(sources))

	for i := range sources {
		ids = append(ids, sources[i].BankID())
	}

	return ids
}
