package pair

import (
	"fmt"

	"github.com/google/uuid"
)

// Pair объединяет prompt и match в одну пару вопроса на соответствие.
type Pair struct {
	prompt Prompt
	match  Match
}

// New создает новую пару и проверяет, что prompt и match инициализированы.
func New(prompt Prompt, match Match) (Pair, error) {
	if err := validatePrompt(prompt); err != nil {
		return Pair{}, fmt.Errorf("ошибка создания pair (prompt): %w", err)
	}

	if err := validateMatch(match); err != nil {
		return Pair{}, fmt.Errorf("ошибка создания pair (match): %w", err)
	}

	return Pair{
		prompt: prompt,
		match:  match,
	}, nil
}

// Prompt возвращает prompt части пары.
func (p Pair) Prompt() Prompt {
	return p.prompt
}

// Match возвращает match части пары.
func (p Pair) Match() Match {
	return p.match
}

// PromptID возвращает идентификатор prompt части пары.
func (p Pair) PromptID() uuid.UUID {
	return p.prompt.ID()
}

// MatchID возвращает идентификатор match части пары.
func (p Pair) MatchID() uuid.UUID {
	return p.match.ID()
}

// IsZero сообщает, что pair находится в нулевом/неинициализированном состоянии.
func (p Pair) IsZero() bool {
	return p.prompt.IsZero() || p.match.IsZero()
}
