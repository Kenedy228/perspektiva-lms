package typed

import (
	"errors"
	"regexp"
	"slices"

	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

const (
	MinPlaceholdersCount int = 1
	MaxPlaceholdersCount int = 20
)

var InTextPlaceholderRegexp = regexp.MustCompile(`\{\{\S+?\}\}`)

var ErrInvalid = errors.New("invalid value")

type Question struct {
	*base.Base
	blanks []blank.Blank
}

func New(t title.Title, blanks []blank.Blank) (*Question, error) {
	placeholders := uniquePlaceholdersInText(t.Value())

	if err := validatePlaceholdersAndBlanks(placeholders, blanks); err != nil {
		return nil, err
	}

	b, err := base.New(t)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:   b,
		blanks: slices.Clone(blanks),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, att *attachment.Attachment, blanks []blank.Blank) (*Question, error) {
	placeholders := uniquePlaceholdersInText(t.Value())

	if err := validatePlaceholdersAndBlanks(placeholders, blanks); err != nil {
		return nil, err
	}

	b, err := base.Restore(id, t, att)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:   b,
		blanks: slices.Clone(blanks),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Blanks() []blank.Blank {
	return slices.Clone(q.blanks)
}

func (q *Question) Type() question2.Type {
	return question2.TypeTyped
}

func (q *Question) ReplaceContent(t title.Title, blanks []blank.Blank) error {
	placeholders := uniquePlaceholdersInText(t.Value())

	if err := validatePlaceholdersAndBlanks(placeholders, blanks); err != nil {
		return err
	}

	q.ChangeTitle(t)
	q.blanks = slices.Clone(blanks)
	return nil
}

func (q *Question) Clone() question2.Question {
	return &Question{
		Base:   q.Base.Clone(),
		blanks: slices.Clone(q.blanks),
	}
}
