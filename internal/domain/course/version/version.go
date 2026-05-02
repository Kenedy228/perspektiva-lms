package version

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/course/version/title"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Version struct {
	id       uuid.UUID
	t        title.Title
	status   Status
	blockIDs []uuid.UUID
}

func New(t title.Title) (*Version, error) {
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

func (v *Version) ChangeTitle(t title.Title) {
	v.t = t
}

func (v *Version) AddBlockID(id uuid.UUID) error {
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

func (v *Version) TryRemoveBlockID(id uuid.UUID) {
	v.blockIDs = slices.DeleteFunc(v.blockIDs, func(current uuid.UUID) bool {
		return current == id
	})
}

func (v *Version) MoveFromTo(from, to int) error {
	if from < 0 || from >= len(v.blockIDs) {
		return fmt.Errorf("%w, детали: неверная начальная позиция перемещаемого блока", ErrInvalid)
	}

	if to < 0 || to >= len(v.blockIDs) {
		return fmt.Errorf("%w, детали: неверная конечная позиция перемещаемого блока", ErrInvalid)
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
		return fmt.Errorf("%w, детали: публиковать можно только черновик", ErrInvalid)
	}

	v.status = StatusPublished
	return nil
}

func (v *Version) MarkAsDeleted() error {
	if v.status == StatusDeleted {
		return fmt.Errorf("%w, детали: версия уже помечена на удаление", ErrInvalid)
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
