package version

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/block"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Version struct {
	id     uuid.UUID
	title  string
	blocks []*block.Block
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

	return &Version{
		id:     id,
		title:  title,
		blocks: cloneBlocks(params.Blocks),
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
	return cloneBlocks(v.blocks)
}

func (v *Version) Status() Status {
	return v.status
}

func (v *Version) ChangeTitle(title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}

	v.title = title
	return nil
}

func (v *Version) InsertBlockAt(pos int, block *block.Block) error {
	if pos < 0 || pos > len(v.blocks) {
		return fmt.Errorf("%w, детали: неверная позиция вставки", ErrInvalid)
	}

	if block == nil {
		return fmt.Errorf("%w, детали: блок для вставки должен существовать", ErrInvalid)
	}

	blocks := slices.Insert(v.blocks, pos, block.Clone())
	if err := validateBlocksDuplication(blocks); err != nil {
		return err
	}

	if err := validateBlocksLimit(blocks, blockLimit); err != nil {
		return err
	}

	v.blocks = blocks
	return nil
}

func (v *Version) RemoveBlockAt(pos int) error {
	if pos < 0 || pos >= len(v.blocks) {
		return fmt.Errorf("%w, детали: неверная позиция удаления", ErrInvalid)
	}

	v.blocks = slices.Delete(v.blocks, pos, pos+1)
	return nil
}

func (v *Version) UpdateBlock(block *block.Block) error {
	for i := range v.blocks {
		if v.blocks[i].ID() == block.ID() {
			v.blocks[i] = block.Clone()
			return nil
		}
	}

	return fmt.Errorf("%w, детали: версия не содержит обновляемый блок", ErrInvalid)
}

func (v *Version) MoveBlock(from, to int) error {
	if v.IsPublished() {
		return fmt.Errorf("%w, детали: вносить изменения в опубликованную версию запрещено", ErrInvalid)
	}

	if from < 0 || from >= len(v.blocks) {
		return fmt.Errorf("%w, детали: неверная начальная позиция перемещаемого элемента", ErrInvalid)
	}

	if to < 0 || to >= len(v.blocks) {
		return fmt.Errorf("%w, детали: неверная конечная позиция перемещаемого элемента", ErrInvalid)
	}

	if from == to {
		return nil
	}

	block := v.blocks[from]
	v.blocks = slices.Delete(v.blocks, from, from+1)
	v.blocks = slices.Insert(v.blocks, to, block)

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

	return &Version{
		id:     id,
		title:  v.title,
		blocks: cloneBlocks(v.blocks),
		status: StatusDraft,
	}, nil
}
