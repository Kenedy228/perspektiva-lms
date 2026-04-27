package attempt

import (
	"time"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Item struct {
	id        uuid.UUID
	snapshot  question.Question
	answer    question.Answer
	changedAt time.Time
	score     int
}

func newItem(question question.Question) (*Item, error) {
	if err := validateSnapshot(question); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}
	snapshot := question.Clone()

	return &Item{
		id:       id,
		snapshot: snapshot,
	}, nil
}

func (i *Item) ID() uuid.UUID {
	return i.id
}

func (i *Item) Snapshot() question.Question {
	return i.snapshot.Clone()
}

func (i *Item) Answer() question.Answer {
	if i.answer == nil {
		return nil
	}
	return i.answer.Clone()
}

func (i *Item) ChangedAt() time.Time {
	return i.changedAt
}

func (i *Item) Score() int {
	return i.score
}

func (i *Item) HasAnswer() bool {
	return i.answer != nil
}

func (i *Item) changeAnswer(answer question.Answer) {
	if answer == nil {
		i.answer = nil
		return
	}

	i.answer = answer.Clone()
	i.touch()
}

func (i *Item) calculateScore() int {
	if i.answer == nil {
		i.score = 0
		return i.score
	}

	check := i.snapshot.CheckAnswer(i.answer)
	if check {
		i.score = 1
	} else {
		i.score = 0
	}

	return i.Score()
}

func (i *Item) touch() {
	i.changedAt = time.Now()
}

func (i *Item) Clone() Item {
	cSnapshot := i.snapshot.Clone()
	var cAnswer question.Answer

	if i.answer == nil {
		cAnswer = nil
	} else {
		cAnswer = i.answer.Clone()
	}

	return Item{
		id:        i.id,
		snapshot:  cSnapshot,
		answer:    cAnswer,
		changedAt: i.changedAt,
		score:     i.score,
	}
}
