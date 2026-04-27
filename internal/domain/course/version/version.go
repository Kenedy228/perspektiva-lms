package version

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/block"
	"gitflic.ru/lms/internal/domain/shared/collection"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Version struct {
	id     uuid.UUID
	title  string
	blocks *collection.OrderedClonable[*block.Block]
	status Status
}

func New(params Params) (*Version, error) {
	title := normalizeTitle(params.Title)

	if err := validateTitle(title); err != nil {
		return nil, err
	}

	if err := validateBlocks(params.Blocks); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	blocks := collection.NewOrderedClonable(params.Blocks)

	return &Version{
		id:     id,
		title:  title,
		blocks: blocks,
		status: StatusDraft,
	}, nil
}

func (v *Version) ID() uuid.UUID {
	return v.id
}

func (v *Version) Title() string {
	return v.title
}

func (v *Version) Blocks() []*block.Block {
	blocks := v.blocks.Items()
	return cloneBlocks(blocks)
}

func (v *Version) Status() Status {
	return v.status
}

func (v *Version) ChangeTitle(title string) error {
	if !v.IsEditable() {
		return ErrNotEditable
	}

	title = normalizeTitle(title)

	if err := validateTitle(title); err != nil {
		return err
	}

	v.title = title
	return nil
}

func (v *Version) InsertBlockAt(block *block.Block, pos int) error {
	if !v.IsEditable() {
		return ErrNotEditable
	}

	if block == nil {
		return fmt.Errorf("%w, детали: блок для вставки должен существовать", ErrInvalid)
	}

	if v.blocks.Count()+1 > blocksLimit {
		return fmt.Errorf("%w, детали: количество блоков в версии не должно превышать %d штук", ErrInvalid, blocksLimit)
	}

	if v.blocks.Contains(block) {
		return fmt.Errorf("%w, детали: версия не должна содержать дубликаты блоков (разрешаются копии)", ErrInvalid)
	}

	err := v.blocks.InsertAt(block, pos)
	if err != nil {
		return fmt.Errorf("%w, детали: %w", ErrInvalid, err)
	}

	return nil
}

func (v *Version) RemoveBlock(block *block.Block) error {
	if !v.IsEditable() {
		return ErrNotEditable
	}

	if block == nil {
		return nil
	}

	err := v.blocks.Remove(block)
	if err != nil {
		return fmt.Errorf("%w, детали: %w", ErrInvalid, err)
	}

	return nil
}

func (v *Version) UpdateBlock(block *block.Block) error {
	if !v.IsEditable() {
		return ErrNotEditable
	}

	if block == nil {
		return fmt.Errorf("%w, детали: блок для обновления должен существовать", ErrInvalid)
	}

	err := v.blocks.Update(block)
	if err != nil {
		return fmt.Errorf("%w, детали: %w", ErrInvalid, err)
	}

	return nil
}

func (v *Version) MoveBlock(from, to int) error {
	if !v.IsEditable() {
		return ErrNotEditable
	}

	err := v.blocks.Move(from, to)
	if err != nil {
		return fmt.Errorf("%w, детали: %w", ErrInvalid, err)
	}

	return nil
}

func (v *Version) IsPublished() bool {
	return v.status == StatusPublished
}

func (v *Version) IsEditable() bool {
	return v.status == StatusDraft
}

func (v *Version) DraftCopy() (*Version, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	blocks := v.blocks.Items()

	return &Version{
		id:     id,
		title:  v.title,
		blocks: collection.NewOrderedClonable(blocks),
		status: StatusDraft,
	}, nil
}
