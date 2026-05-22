package matching

import (
	"slices"

	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type Question struct {
	*base.Base
	pairs []pair.Pair
}

func New(t title.Title, pairs []pair.Pair) (*Question, error) {
	if err := validatePairs(pairs); err != nil {
		return nil, err
	}

	b, err := base.New(t)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:  b,
		pairs: slices.Clone(pairs),
	}, nil
}

func Restore(id uuid.UUID, t title.Title, att *attachment.Attachment, pairs []pair.Pair) (*Question, error) {
	if err := validatePairs(pairs); err != nil {
		return nil, err
	}

	b, err := base.Restore(id, t, att)
	if err != nil {
		return nil, err
	}

	return &Question{
		Base:  b,
		pairs: slices.Clone(pairs),
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Pairs() []pair.Pair {
	return slices.Clone(q.pairs)
}

func (q *Question) Type() question2.Type {
	return question2.TypeMatching
}

func (q *Question) ChangePairs(pairs []pair.Pair) error {
	if err := validatePairs(pairs); err != nil {
		return err
	}

	q.pairs = slices.Clone(pairs)
	return nil
}

func (q *Question) Clone() question2.Question {
	return &Question{
		Base:  q.Base.Clone(),
		pairs: slices.Clone(q.pairs),
	}
}
