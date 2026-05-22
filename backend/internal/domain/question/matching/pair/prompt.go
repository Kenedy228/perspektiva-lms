package pair

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Prompt struct {
	id   uuid.UUID
	text text.Text
}

const (
	PromptCharsLimit int = 255
)

func NewPrompt(t text.Text) (Prompt, error) {
	if utf8.RuneCountInString(t.Value()) > PromptCharsLimit {
		return Prompt{}, fmt.Errorf("%w: invalid value (%d)", ErrInvalid, PromptCharsLimit)
	}

	id, err := uid.New()
	if err != nil {
		return Prompt{}, err
	}

	return Prompt{
		id:   id,
		text: t,
	}, nil
}

func RestorePrompt(id uuid.UUID, t text.Text) (Prompt, error) {
	if id == uuid.Nil {
		return Prompt{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if t.IsIncomplete() {
		return Prompt{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if utf8.RuneCountInString(t.Value()) > PromptCharsLimit {
		return Prompt{}, fmt.Errorf("%w: invalid value (%d)", ErrInvalid, PromptCharsLimit)
	}

	return Prompt{
		id:   id,
		text: t,
	}, nil
}

func (p Prompt) ID() uuid.UUID {
	return p.id
}

func (p Prompt) Text() text.Text {
	return p.text
}

func (p Prompt) IsIncomplete() bool {
	return p.id == uuid.Nil || len(p.text.Value()) == 0
}
