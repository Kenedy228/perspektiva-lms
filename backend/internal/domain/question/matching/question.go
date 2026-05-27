package matching

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
)

type Question struct {
	*base.Base
	pairs []pair.Pair
}

// New создает новый вопрос на соответствие и валидирует базовые данные и пары.
func New(b *base.Base, pairs []pair.Pair) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validatePairs(pairs); err != nil {
		return nil, err
	}

	return &Question{
		Base:  b,
		pairs: slices.Clone(pairs),
	}, nil
}

// Restore восстанавливает вопрос на соответствие из существующего состояния.
func Restore(b *base.Base, pairs []pair.Pair) (*Question, error) {
	if err := validateBase(b); err != nil {
		return nil, err
	}

	if err := validatePairs(pairs); err != nil {
		return nil, err
	}

	return &Question{
		Base:  b,
		pairs: slices.Clone(pairs),
	}, nil
}

// Instruction возвращает стандартную инструкцию для вопроса matching.
func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

// Pairs возвращает копию пар вопроса.
func (q *Question) Pairs() []pair.Pair {
	return slices.Clone(q.pairs)
}

// Type возвращает тип вопроса.
func (q *Question) Type() question.Type {
	return question.TypeMatching
}

// ChangePairs заменяет пары вопроса после валидации.
func (q *Question) ChangePairs(pairs []pair.Pair) error {
	if err := validatePairs(pairs); err != nil {
		return err
	}

	q.pairs = slices.Clone(pairs)
	return nil
}

// Clone создает полную копию вопроса.
func (q *Question) Clone() question.Question {
	return &Question{
		Base:  q.Base.Clone(),
		pairs: slices.Clone(q.pairs),
	}
}
