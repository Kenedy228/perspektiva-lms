package element

import (
	"gitflic.ru/lms/backend/internal/domain/course/element/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Element struct {
	id             uuid.UUID
	t              title.Title
	content        Content
	completionMode CompletionMode
}

func New(t title.Title, c Content) (*Element, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateContent(c); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Element{
		id:             id,
		t:              t,
		content:        c.Clone(),
		completionMode: CompletionModeNone,
	}, nil
}

func Restore(id uuid.UUID, t title.Title, c Content) (*Element, error) {
	return RestoreWithCompletionMode(id, t, c, CompletionModeNone)
}

func RestoreWithCompletionMode(id uuid.UUID, t title.Title, c Content, mode CompletionMode) (*Element, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateContent(c); err != nil {
		return nil, err
	}
	if err := validateCompletionMode(mode); err != nil {
		return nil, err
	}

	return &Element{
		id:             id,
		t:              t,
		content:        c.Clone(),
		completionMode: mode,
	}, nil
}

func (e *Element) ID() uuid.UUID {
	return e.id
}

func (e *Element) Title() title.Title {
	return e.t
}

func (e *Element) Content() Content {
	if e.content == nil {
		return nil
	}
	return e.content.Clone()
}

func (e *Element) ChangeTitle(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}
	e.t = t
	return nil
}

func (e *Element) ChangeContent(c Content) error {
	if err := validateContent(c); err != nil {
		return err
	}

	e.content = c.Clone()
	return nil
}

func (e *Element) CompletionMode() CompletionMode {
	return e.completionMode
}

func (e *Element) ChangeCompletionMode(mode CompletionMode) error {
	if err := validateCompletionMode(mode); err != nil {
		return err
	}
	e.completionMode = mode
	return nil
}

func (e *Element) IsTrackable() bool {
	return e.completionMode == CompletionModeManual
}

func (e *Element) Replicate() (*Element, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Element{
		id:             id,
		t:              e.t,
		content:        e.content.Clone(),
		completionMode: e.completionMode,
	}, nil
}

func (e *Element) Clone() *Element {
	return &Element{
		id:             e.id,
		t:              e.t,
		content:        e.content.Clone(),
		completionMode: e.completionMode,
	}
}

func (e *Element) IsInteractive() bool {
	if e.content == nil {
		return false
	}
	return e.content.IsInteractive()
}
