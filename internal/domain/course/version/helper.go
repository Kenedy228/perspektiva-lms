package version

import "gitflic.ru/lms/internal/domain/block"

func cloneBlocks(original []*block.Block) []*block.Block {
	blocks := make([]*block.Block, 0, len(original))

	for i := range original {
		blocks = append(blocks, original[i])
	}

	return blocks
}
