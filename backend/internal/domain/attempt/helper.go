package attempt

import (
	"gitflic.ru/lms/backend/internal/domain/attempt/item"
	"github.com/google/uuid"
)

func findItem(items []item.Item, questionID uuid.UUID) (item.Item, bool) {
	for i := range items {
		if items[i].ID() == questionID {
			return items[i], true
		}
	}
	return item.Item{}, false
}
