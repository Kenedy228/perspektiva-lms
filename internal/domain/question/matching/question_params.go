package matching

import (
	"gitflic.ru/lms/internal/domain/question/base"
	"github.com/google/uuid"
)

type Params struct {
	Text    string
	ImageID uuid.UUID
	Pairs   []Pair
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:    p.Text,
		ImageID: p.ImageID,
	}
}
