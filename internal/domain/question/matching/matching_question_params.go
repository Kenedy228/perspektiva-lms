package matching

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/google/uuid"
)

type Params struct {
	Text  question.QText
	Image uuid.UUID
	Pairs []PairParams
}

func (p Params) baseParams() base.Params {
	return base.Params{
		Text:        p.Text,
		Description: question.QDescriptionMatching,
		ImageID:     p.Image,
	}
}

type PairParams struct {
	Prompt        string
	ContentOption option.ContentOption
}
