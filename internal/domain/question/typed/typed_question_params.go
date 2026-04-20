package typed

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/google/uuid"
)

type Params struct {
	Text    question.QText
	ImageID uuid.UUID
	Blanks  []BlankParams
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:        p.Text,
		Description: question.QDescriptionTyped,
		ImageID:     p.ImageID,
	}
}

type BlankParams struct {
	Placeholder string
	Answers     []option.ContentOption
}
