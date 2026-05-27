package pair

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

// Match хранит правую часть пары вопроса на соответствие.
type Match struct {
	id    uuid.UUID
	value string
}

// NewMatch создает новый match и проверяет его значение.
func NewMatch(value string) (Match, error) {
	if err := validateValue(value); err != nil {
		return Match{}, fmt.Errorf("ошибка создания match: %w", err)
	}

	id, err := uid.New()
	if err != nil {
		return Match{}, err
	}

	return Match{
		id:    id,
		value: value,
	}, nil
}

// RestoreMatch восстанавливает match из сохраненного состояния.
func RestoreMatch(id uuid.UUID, value string) (Match, error) {
	if err := validateIDRequired(id); err != nil {
		return Match{}, fmt.Errorf("ошибка восстановления match: %w", err)
	}

	if err := validateValue(value); err != nil {
		return Match{}, fmt.Errorf("ошибка восстановления match: %w", err)
	}

	return Match{
		id:    id,
		value: value,
	}, nil
}

// ID возвращает идентификатор match.
func (m Match) ID() uuid.UUID {
	return m.id
}

// Value возвращает значение match.
func (m Match) Value() string {
	return m.value
}

// IsZero сообщает, что match находится в нулевом/неинициализированном состоянии.
func (m Match) IsZero() bool {
	return m.id == uuid.Nil || len(m.value) == 0
}
