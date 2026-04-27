package version

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/block"
	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/google/uuid"
)

func validateBlocks(blocks []*block.Block) error {
	if err := validateBlocksLimit(blocks); err != nil {
		return err
	}

	if err := validateBlocksNil(blocks); err != nil {
		return err
	}

	if err := validateBlocksDuplication(blocks); err != nil {
		return err
	}

	return nil
}

func validateBlocksDuplication(blocks []*block.Block) error {
	ids := make([]uuid.UUID, 0, len(blocks))

	for i := range blocks {
		ids = append(ids, blocks[i].ID())
	}

	if has := duplicate.FindUUID(ids); has {
		return fmt.Errorf("%w, детали: версия не должна содержать дубликаты блоков (разрешаются копии)", ErrInvalid)
	}

	return nil
}

func validateBlocksNil(blocks []*block.Block) error {
	hasNil := slices.ContainsFunc(blocks, func(current *block.Block) bool {
		return current == nil
	})

	if hasNil {
		return fmt.Errorf("%w, детали: версия не должна содержать не существующие блоки", ErrInvalid)
	}

	return nil
}

func validateBlocksLimit(blocks []*block.Block) error {
	if len(blocks) > blocksLimit {
		return fmt.Errorf("%w, детали: количество блоков в версии не должно превышать %d штук", ErrInvalid, blocksLimit)
	}

	return nil
}
