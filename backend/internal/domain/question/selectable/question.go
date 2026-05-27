package selectable

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
)

type Question struct {
	*base.Base
	options []option.Option
}

func New(b *base.Base, options []option.Option) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validateOptions(options); err != nil {
		return nil, err
	}

	return &Question{
		Base:    b,
		options: slices.Clone(options),
	}, nil
}

func Restore(b *base.Base, options []option.Option) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validateOptions(options); err != nil {
		return nil, err
	}

	return &Question{
		Base:    b,
		options: slices.Clone(options),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Options() []option.Option {
	return slices.Clone(q.options)
}

func (q *Question) CorrectOptionsCount() int {
	return countCorrectOptions(q.options)
}

func (q *Question) Type() question.Type {
	return question.TypeSelectable
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
