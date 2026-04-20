package typed

import (
	"slices"
	"strings"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
)

const (
	minPlaceholders          = 2
	maxPlaceholders          = 20
	maxAnswersPerPlaceholder = 20
)

type TypedQuestion struct {
	base.Base
	blanks []Blank
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validatePlaceholders(params.Text.String(), params.Blanks); err != nil {
		return nil, err
	}

	blanks, err := mapBlanks(params.Blanks)
	if err != nil {
		return nil, err
	}

	return &TypedQuestion{
		Base:   base,
		blanks: blanks,
	}, nil
}

func (q *TypedQuestion) Blanks() []Blank {
	return slices.Clone(q.blanks)
}

func (q *TypedQuestion) ReplaceContent(text question.QText, rawBlanks []BlankParams) error {
	if err := validatePlaceholders(text.String(), rawBlanks); err != nil {
		return err
	}

	blanks, err := mapBlanks(rawBlanks)
	if err != nil {
		return err
	}

	q.UpdateText(text)
	q.blanks = blanks

	return nil
}

func (q *TypedQuestion) Type() question.Type {
	return question.TypeTyped
}

func (q *TypedQuestion) HasAnswerForPlaceholder(placeholder string, answer option.ContentOption) bool {
	for i := range q.blanks {
		if q.blanks[i].Placeholder() == placeholder {
			return slices.ContainsFunc(q.blanks[i].answers, func(current option.ContentOption) bool {
				return strings.EqualFold(current.Value(), answer.Value())
			})
		}
	}

	return false
}
