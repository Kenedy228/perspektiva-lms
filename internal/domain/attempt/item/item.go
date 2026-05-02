package item

import (
	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Item struct {
	snapshot question.Question
}

func New(q question.Question) (Item, error) {
	if err := validateQuestion(q); err != nil {
		return Item{}, err
	}

	return Item{
		snapshot: q.Clone(),
	}, nil
}

func (i Item) ID() uuid.UUID {
	return i.snapshot.ID()
}

func (i Item) Snapshot() question.Question {
	return i.snapshot.Clone()
}
