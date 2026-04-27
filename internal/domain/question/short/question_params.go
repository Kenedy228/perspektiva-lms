package short

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"github.com/google/uuid"
)

type Params struct {
	Text     string
	ImageID  uuid.UUID
	Variants []question.Content
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:    p.Text,
		ImageID: p.ImageID,
	}
}
