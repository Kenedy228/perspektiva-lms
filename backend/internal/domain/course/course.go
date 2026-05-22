package course

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Course struct {
	id         uuid.UUID
	t          title.Title
	versionIDs []uuid.UUID
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
		id:         id,
		t:          t,
		versionIDs: make([]uuid.UUID, 0),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, versionIDs []uuid.UUID) (*Course, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateVersionIDs(versionIDs); err != nil {
		return nil, err
	}

	return &Course{
		id:         id,
		t:          t,
		versionIDs: slices.Clone(versionIDs),
	}, nil
}

func (c *Course) ID() uuid.UUID {
	return c.id
}

func (c *Course) Title() title.Title {
	return c.t
}

func (c *Course) VersionIDs() []uuid.UUID {
	return slices.Clone(c.versionIDs)
}

func (c *Course) ChangeTitle(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}
	c.t = t
	return nil
}

func (c *Course) AddVersionID(versionID uuid.UUID) error {
	if err := validateRequiredVersionID(versionID); err != nil {
		return err
	}

	if err := validateVersionIDsLimit(c.versionIDs); err != nil {
		return err
	}

	if err := validateVersionIDsDuplication(versionID, c.versionIDs); err != nil {
		return err
	}

	c.versionIDs = append(c.versionIDs, versionID)
	return nil
}

func (c *Course) RemoveVersionID(versionID uuid.UUID) error {
	if err := validateRequiredVersionID(versionID); err != nil {
		return err
	}
	if !slices.Contains(c.versionIDs, versionID) {
		return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, versionID)
	}
	c.versionIDs = slices.DeleteFunc(c.versionIDs, func(current uuid.UUID) bool {
		return current == versionID
	})
	return nil
}

func (c *Course) HasVersion(versionID uuid.UUID) bool {
	return slices.Contains(c.versionIDs, versionID)
}

func (c *Course) Clone() *Course {
	return &Course{
		id:         c.id,
		t:          c.t,
		versionIDs: slices.Clone(c.versionIDs),
	}
}
