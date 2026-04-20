package short

import (
	"slices"
	"strings"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
)

const maxAnswers = 20

type ShortQuestion struct {
	base.Base
	answers []option.ContentOption
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateAnswers(params.Answers); err != nil {
		return nil, err
	}

	cAnswers := slices.Clone(params.Answers)

	return &ShortQuestion{
		Base:    base,
		answers: cAnswers,
	}, nil
}

func (q *ShortQuestion) Answers() []option.ContentOption {
	return slices.Clone(q.answers)
}

func (q *ShortQuestion) UpdateAnswers(answers []option.ContentOption) error {
	if err := validateAnswers(answers); err != nil {
		return err
	}

	cAnswers := slices.Clone(answers)
	q.answers = cAnswers
	q.Touch()
	return nil
}

func (q *ShortQuestion) Type() question.Type {
	return question.TypeShort
}

func (q *ShortQuestion) HasAnswer(answer option.ContentOption) bool {
	if answer.ContentType() != option.ContentTypeText {
		return false
	}

	return slices.ContainsFunc(q.answers, func(current option.ContentOption) bool {
		return strings.EqualFold(current.Value(), answer.Value())
	})
}
