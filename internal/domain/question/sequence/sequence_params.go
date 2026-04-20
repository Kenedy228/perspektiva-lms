package sequence

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/google/uuid"
)

type Params struct {
	Text    question.QText
	ImageID uuid.UUID
	Items   []option.ContentOption
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:        p.Text,
		Description: question.QDescriptionSequence,
		ImageID:     p.ImageID,
	}
}
