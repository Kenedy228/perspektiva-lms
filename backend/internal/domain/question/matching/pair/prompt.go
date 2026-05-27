package pair

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

// Prompt хранит левую часть пары вопроса на соответствие.
type Prompt struct {
	id    uuid.UUID
	value string
}

// NewPrompt создает новый prompt и проверяет его значение.
func NewPrompt(value string) (Prompt, error) {
	if err := validateValue(value); err != nil {
		return Prompt{}, fmt.Errorf("ошибка создания prompt: %w", err)
	}

	id, err := uid.New()
	if err != nil {
		return Prompt{}, err
	}

	return Prompt{
		id:    id,
		value: value,
	}, nil
}

// RestorePrompt восстанавливает prompt из сохраненного состояния.
func RestorePrompt(id uuid.UUID, value string) (Prompt, error) {
	if err := validateIDRequired(id); err != nil {
		return Prompt{}, fmt.Errorf("ошибка восстановления prompt: %w", err)
	}

	if err := validateValue(value); err != nil {
		return Prompt{}, fmt.Errorf("ошибка восстановления prompt: %w", err)
	}

	return Prompt{
		id:    id,
		value: value,
	}, nil
}

// ID возвращает идентификатор prompt.
func (p Prompt) ID() uuid.UUID {
	return p.id
}

// Value возвращает значение prompt.
func (p Prompt) Value() string {
	return p.value
}

// IsZero сообщает, что prompt находится в нулевом/неинициализированном состоянии.
func (p Prompt) IsZero() bool {
	return p.id == uuid.Nil || len(p.value) == 0
}
