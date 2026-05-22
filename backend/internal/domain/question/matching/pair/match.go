package pair

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Match struct {
	id   uuid.UUID
	text text.Text
}

const (
	MatchCharsLimit int = 255
)

func NewMatch(t text.Text) (Match, error) {
	if utf8.RuneCountInString(t.Value()) > MatchCharsLimit {
		return Match{}, fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MatchCharsLimit)
	}

	id, err := uid.New()
	if err != nil {
		return Match{}, err
	}

	return Match{
		id:   id,
		text: t,
	}, nil
}

func RestoreMatch(id uuid.UUID, t text.Text) (Match, error) {
	if id == uuid.Nil {
		return Match{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if t.IsIncomplete() {
		return Match{}, fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if utf8.RuneCountInString(t.Value()) > MatchCharsLimit {
		return Match{}, fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MatchCharsLimit)
	}

	return Match{
		id:   id,
		text: t,
	}, nil
}

func (m Match) ID() uuid.UUID {
	return m.id
}

func (m Match) Text() text.Text {
	return m.text
}

func (m Match) IsIncomplete() bool {
	return m.id == uuid.Nil || len(m.text.Value()) == 0
}
