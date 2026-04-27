package version

import "gitflic.ru/lms/internal/domain/block"

type Params struct {
	Title  string
	Blocks []*block.Block
}
