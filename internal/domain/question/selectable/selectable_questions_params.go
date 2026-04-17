package selectable

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/google/uuid"
)

type ItemParams struct {
	Content    option.ContentOption
	IsCorrect bool
}

type Params struct {
	Text    question.QText
	ImageID uuid.UUID
	Items   []ItemParams
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:        p.Text,
		Description: question.QDescriptionSelectable,
		ImageID:     p.ImageID,
	}
}
