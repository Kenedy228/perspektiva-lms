package sequence

import (
	"slices"

	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type Question struct {
	*base.Base
	options []option.Option
}

func New(t title.Title, options []option.Option) (*Question, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}

	b, err := base.New(t)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:    b,
		options: slices.Clone(options),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, att *attachment.Attachment, options []option.Option) (*Question, error) {
	if err := validateOptions(options); err != nil {
		return nil, err
	}

	b, err := base.Restore(id, t, att)
	if err != nil {
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

func (q *Question) Type() question2.Type {
	return question2.TypeSequence
}

func (q *Question) ChangeOptions(options []option.Option) error {
	if err := validateOptions(options); err != nil {
		return err
	}

	q.options = slices.Clone(options)
	return nil
}

func (q *Question) Clone() question2.Question {
	return &Question{
		Base:    q.Base.Clone(),
		options: slices.Clone(q.options),
	}
}
