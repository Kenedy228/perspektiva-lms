package block

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/block/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Block struct {
	id         uuid.UUID
	t          title.Title
	elementIDs []uuid.UUID
}

func New(t title.Title) (*Block, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Block{
		id:         id,
		t:          t,
		elementIDs: make([]uuid.UUID, 0),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, elementIDs []uuid.UUID) (*Block, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateElementIDs(elementIDs); err != nil {
		return nil, err
	}

	return &Block{
		id:         id,
		t:          t,
		elementIDs: slices.Clone(elementIDs),
	}, nil
}

func (b *Block) ID() uuid.UUID {
	return b.id
}

func (b *Block) Title() title.Title {
	return b.t
}

func (b *Block) ElementIDs() []uuid.UUID {
	return slices.Clone(b.elementIDs)
}

func (b *Block) ChangeTitle(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}
	b.t = t
	return nil
}

func (b *Block) AddElementID(id uuid.UUID) error {
	if err := validateRequiredElementID(id); err != nil {
		return err
	}

	if err := validateElementIDsLimit(b.elementIDs); err != nil {
		return err
	}

	if err := validateElementIDsDuplication(id, b.elementIDs); err != nil {
		return err
	}

	b.elementIDs = append(b.elementIDs, id)
	return nil
}

func (b *Block) RemoveElementID(id uuid.UUID) error {
	if err := validateRequiredElementID(id); err != nil {
		return err
	}
	if !slices.Contains(b.elementIDs, id) {
		return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, id)
	}
	b.elementIDs = slices.DeleteFunc(b.elementIDs, func(current uuid.UUID) bool {
		return current == id
	})
	return nil
}

func (b *Block) MoveFromTo(from, to int) error {
	if from < 0 || from >= len(b.elementIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if to < 0 || to >= len(b.elementIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if from == to {
		return nil
	}

	id := b.elementIDs[from]
	b.elementIDs = slices.Delete(b.elementIDs, from, from+1)
	b.elementIDs = slices.Insert(b.elementIDs, to, id)
	return nil
}

func (b *Block) Clone() *Block {
	return &Block{
		id:         b.id,
		t:          b.t,
		elementIDs: slices.Clone(b.elementIDs),
	}
}

func (b *Block) Replicate() (*Block, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Block{
		id:         id,
		t:          b.t,
		elementIDs: slices.Clone(b.elementIDs),
	}, nil
}
