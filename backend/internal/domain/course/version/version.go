package version

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/course/version/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Version struct {
	id       uuid.UUID
	t        title.Title
	status   Status
	blockIDs []uuid.UUID
}

func New(t title.Title) (*Version, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Version{
		id:       id,
		t:        t,
		blockIDs: make([]uuid.UUID, 0),
		status:   StatusDraft,
	}, nil
}

func Restore(id uuid.UUID, t title.Title, status Status, blockIDs []uuid.UUID) (*Version, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	if err := validateTitle(t); err != nil {
		return nil, err
	}
	if err := validateStatus(status); err != nil {
		return nil, err
	}
	if err := validateBlockIDs(blockIDs); err != nil {
		return nil, err
	}

	return &Version{
		id:       id,
		t:        t,
		status:   status,
		blockIDs: slices.Clone(blockIDs),
	}, nil
}

func (v *Version) ID() uuid.UUID {
	return v.id
}

func (v *Version) Title() title.Title {
	return v.t
}

func (v *Version) BlockIDs() []uuid.UUID {
	return slices.Clone(v.blockIDs)
}

func (v *Version) Status() Status {
	return v.status
}

func (v *Version) ChangeTitle(t title.Title) error {
	if err := validateDraft(v.status); err != nil {
		return err
	}
	if err := validateTitle(t); err != nil {
		return err
	}
	v.t = t
	return nil
}

func (v *Version) AddBlockID(id uuid.UUID) error {
	if err := validateDraft(v.status); err != nil {
		return err
	}

	if err := validateRequiredBlockID(id); err != nil {
		return err
	}

	if err := validateBlockIDsLimit(v.blockIDs); err != nil {
		return err
	}

	if err := validateBlockIDsDuplication(id, v.blockIDs); err != nil {
		return err
	}

	v.blockIDs = append(v.blockIDs, id)
	return nil
}

func (v *Version) RemoveBlockID(id uuid.UUID) error {
	if err := validateDraft(v.status); err != nil {
		return err
	}
	if err := validateRequiredBlockID(id); err != nil {
		return err
	}
	if !slices.Contains(v.blockIDs, id) {
		return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, id)
	}
	v.blockIDs = slices.DeleteFunc(v.blockIDs, func(current uuid.UUID) bool {
		return current == id
	})
	return nil
}

func (v *Version) MoveFromTo(from, to int) error {
	if err := validateDraft(v.status); err != nil {
		return err
	}

	if from < 0 || from >= len(v.blockIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if to < 0 || to >= len(v.blockIDs) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if from == to {
		return nil
	}

	id := v.blockIDs[from]
	v.blockIDs = slices.Delete(v.blockIDs, from, from+1)
	v.blockIDs = slices.Insert(v.blockIDs, to, id)
	return nil
}

func (v *Version) Publish() error {
	if v.status != StatusDraft {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	if len(v.blockIDs) == 0 {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	v.status = StatusPublished
	return nil
}

func (v *Version) MarkAsDeleted() error {
	if v.status == StatusDeleted {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	v.status = StatusDeleted
	return nil
}

func (v *Version) IsDraft() bool {
	return v.status == StatusDraft
}

func (v *Version) IsPublished() bool {
	return v.status == StatusPublished
}

func (v *Version) IsDeleted() bool {
	return v.status == StatusDeleted
}

func (v *Version) Clone() *Version {
	return &Version{
		id:       v.id,
		t:        v.t,
		blockIDs: slices.Clone(v.blockIDs),
		status:   v.status,
	}
}

func (v *Version) Replicate() (*Version, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Version{
		id:       id,
		t:        v.t,
		blockIDs: slices.Clone(v.blockIDs),
		status:   StatusDraft,
	}, nil
}
