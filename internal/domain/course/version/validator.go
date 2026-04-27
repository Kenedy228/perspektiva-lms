package version

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/block"
	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/google/uuid"
)

func validateTitle(title string) error {
	if err := validateNotEmptyStr(title); err != nil {
		return err
	}

	if err := validateLimitStr(title, titleCharsLimit); err != nil {
		return err
	}

	return nil
}

func validateBlocks(blocks []*block.Block) error {
	if err := validateBlocksLimit(blocks, blockLimit); err != nil {
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

func validateNotEmptyStr(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w, детали: заголовок блока должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateLimitStr(s string, limit int) error {
	if utf8.RuneCountInString(s) > limit {
		return fmt.Errorf("%w, детали: заголовок должен содержать не более %d символов", ErrInvalid, limit)
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
	for i := range blocks {
		if blocks[i] == nil {
			return fmt.Errorf("%w, детали: версия не должна содержать не существующие блоки", ErrInvalid)
		}
	}

	return nil
}

func validateBlocksLimit(blocks []*block.Block, limit int) error {
	if len(blocks) > limit {
		return fmt.Errorf("%w, детали: количество блоков в версии не должно превышать %d штук", ErrInvalid, limit)
	}

	return nil
}
