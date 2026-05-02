package course

import (
	"slices"

	"gitflic.ru/lms/internal/domain/course/title"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Course struct {
	id         uuid.UUID
	t          title.Title
	versionIDs []uuid.UUID
}

func New(t title.Title) (*Course, error) {
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

func (c *Course) ID() uuid.UUID {
	return c.id
}

func (c *Course) Title() title.Title {
	return c.t
}

func (c *Course) VersionIDs() []uuid.UUID {
	return slices.Clone(c.versionIDs)
}

func (c *Course) ChangeTitle(t title.Title) {
	c.t = t
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

func (c *Course) TryRemoveVersionID(versionID uuid.UUID) {
	c.versionIDs = slices.DeleteFunc(c.versionIDs, func(current uuid.UUID) bool {
		return current == versionID
	})
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
