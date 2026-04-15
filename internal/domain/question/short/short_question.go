package short

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const description = "введите короткий ответ"
const maxAnswers = 20

type ShortQuestion struct {
	base.Base
	answers []string
}

func New(params *Params) (question.Question, error) {
	if err := validateAnswers(params.Answers, params.AllowDuplicate); err != nil {
		return nil, err
	}

	base, err := base.New(&base.Params{
		Text:        params.Text,
		Description: description,
		Image:       params.Image,
	})

	if err != nil {
		return nil, err
	}

	answersCopy := slices.Clone(params.Answers)

	return &ShortQuestion{
		Base:    base,
		answers: answersCopy,
	}, nil
}

func (q *ShortQuestion) Answers() []string {
	return slices.Clone(q.answers)
}

func (q *ShortQuestion) UpdateAnswers(answers []string, allowDuplicate bool) error {
	if err := validateAnswers(answers, allowDuplicate); err != nil {
		return err
	}

	q.answers = slices.Clone(answers)
	q.Touch()
	return nil
}

func (q *ShortQuestion) Type() question.Type {
	return question.TypeShort
}
