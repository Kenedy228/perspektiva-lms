package sequence

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/sequence/option"
	"gitflic.ru/lms/internal/domain/question/title"
)

type Question struct {
	*base.Base
	options []option.Option
}

func New(t title.Title, options []option.Option) (*Question, error) {
	base, err := base.New(t)
	if err != nil {
		return nil, err
	}

	if err := validateOptions(options); err != nil {
		return nil, err
	}

	return &Question{
		Base:    base,
		options: slices.Clone(options),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Options() []option.Option {
	return slices.Clone(q.options)
}

func (q *Question) Type() question.Type {
	return question.TypeSequence
}

func (q *Question) ChangeOptions(options []option.Option) error {
	if err := validateOptions(options); err != nil {
		return err
	}

	q.options = slices.Clone(options)
	return nil
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:    q.Base.Clone(),
		options: slices.Clone(q.options),
	}
}
