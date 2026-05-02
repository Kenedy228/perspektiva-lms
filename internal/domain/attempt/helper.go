package attempt

import (
	"gitflic.ru/lms/internal/domain/attempt/item"
	"github.com/google/uuid"
)

func hasQuestion(items []item.Item, questionID uuid.UUID) bool {
	for i := range items {
		if items[i].ID() == questionID {
			return true
		}
	}
	return false
}
