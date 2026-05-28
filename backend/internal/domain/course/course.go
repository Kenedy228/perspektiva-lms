package course

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Course struct {
	id       uuid.UUID
	t        title.Title
	blockIDs []uuid.UUID
}

func New(t title.Title) (*Course, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Course{
		id:       id,
		t:        t,
		blockIDs: make([]uuid.UUID, 0),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, blockIDs []uuid.UUID) (*Course, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateBlockIDs(blockIDs); err != nil {
		return nil, err
	}

	return &Course{
		id:       id,
		t:        t,
		blockIDs: slices.Clone(blockIDs),
	}, nil
}

func (c *Course) ID() uuid.UUID {
	return c.id
}

func (c *Course) Title() title.Title {
	return c.t
}

func (c *Course) BlockIDs() []uuid.UUID {
	return slices.Clone(c.blockIDs)
}

func (c *Course) ChangeTitle(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}
	c.t = t
	return nil
}

func (c *Course) AddBlockID(blockID uuid.UUID) error {
	if err := validateRequiredBlockID(blockID); err != nil {
		return err
	}

	if err := validateBlockIDsLimit(c.blockIDs); err != nil {
		return err
	}

	if err := validateBlockIDsDuplication(blockID, c.blockIDs); err != nil {
		return err
	}

	c.blockIDs = append(c.blockIDs, blockID)
	return nil
}

func (c *Course) RemoveBlockID(blockID uuid.UUID) error {
	if err := validateRequiredBlockID(blockID); err != nil {
		return err
	}
	if !slices.Contains(c.blockIDs, blockID) {
		return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, blockID)
	}
	c.blockIDs = slices.DeleteFunc(c.blockIDs, func(current uuid.UUID) bool {
		return current == blockID
	})
	return nil
}

func (c *Course) MoveBlock(from, to int) error {
	if from < 0 || from >= len(c.blockIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	if to < 0 || to >= len(c.blockIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	if from == to {
		return nil
	}

	id := c.blockIDs[from]
	c.blockIDs = slices.Delete(c.blockIDs, from, from+1)
	c.blockIDs = slices.Insert(c.blockIDs, to, id)
	return nil
}

func (c *Course) Clone() *Course {
	return &Course{
		id:       c.id,
		t:        c.t,
		blockIDs: slices.Clone(c.blockIDs),
	}
}
