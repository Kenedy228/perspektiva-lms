// Пакет uid предоставляет точку входа для генерации UUIDv7 для всех сущностей в приложении.
package uid

import (
	"fmt"

	"github.com/google/uuid"
)

// New генерирует UUIDv7 и возвращает сгенерированное значение.
// Если при генерации возникла ошибка, то возвращается значение uuid.Nil.
func New() (uuid.UUID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("ошибка при генерации идентификатора, детали: %w", err)
	}

	return id, nil
}
