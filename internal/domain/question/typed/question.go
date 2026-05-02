package typed

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/title"
	"gitflic.ru/lms/internal/domain/question/typed/blank"
)

type Question struct {
	*base.Base
	blanks []blank.Blank
}

func New(t title.Title, blanks []blank.Blank) (*Question, error) {
	base, err := base.New(t)
	if err != nil {
		return nil, err
	}

	placeholders := findPlaceholdersInText(t.Value())

	if err := validatePlaceholders(placeholders, blanks); err != nil {
		return nil, err
	}

	return &Question{
		Base:   base,
		blanks: slices.Clone(blanks),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Blanks() []blank.Blank {
	return slices.Clone(q.blanks)
}

func (q *Question) Type() question.Type {
	return question.TypeTyped
}

func (q *Question) ReplaceContent(t title.Title, blanks []blank.Blank) error {
	placeholders := findPlaceholdersInText(t.Value())
	if err := validatePlaceholders(placeholders, blanks); err != nil {
		return err
	}

	q.ChangeTitle(t)
	q.blanks = slices.Clone(blanks)
	return nil
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:   q.Base.Clone(),
		blanks: slices.Clone(q.blanks),
	}
}
