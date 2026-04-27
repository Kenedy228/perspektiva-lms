package typed

import (
	"slices"
	"strings"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

type Question struct {
	*base.Base
	blanks []Blank
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validatePlaceholders(params.Text, params.Blanks); err != nil {
		return nil, err
	}

	cBlanks := slices.Clone(params.Blanks)

	return &Question{
		Base:   base,
		blanks: cBlanks,
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Blanks() []Blank {
	return slices.Clone(q.blanks)
}

func (q *Question) ReplaceContent(text string, blanks []Blank) error {
	if err := validatePlaceholders(text, blanks); err != nil {
		return err
	}

	cBlanks := slices.Clone(blanks)
	q.UpdateText(text)
	q.blanks = cBlanks

	return nil
}

func (q *Question) Type() question.Type {
	return question.TypeTyped
}

func (q *Question) CheckAnswer(answer question.Answer) bool {
	cast, ok := answer.(Answer)
	if !ok {
		return false
	}

	studentAnswers := cast.Inputs()

	if len(studentAnswers) != len(q.blanks) {
		return false
	}

	for i := range q.blanks {
		v, ok := studentAnswers[q.blanks[i].Placeholder()]

		if !ok {
			return false
		}

		isEqual := slices.ContainsFunc(q.blanks[i].Variants(), func(opt question.Content) bool {
			return strings.EqualFold(opt.Value(), v)
		})

		if !isEqual {
			return false
		}
	}

	return true
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:   q.Base.Clone(),
		blanks: q.Blanks(),
	}
}
