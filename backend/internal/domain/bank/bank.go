package bank

import (
	"slices"

	"gitflic.ru/lms/backend/internal/domain/bank/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Bank struct {
	id        uuid.UUID
	t         title.Title
	questions []uuid.UUID
}

func New(t title.Title) (*Bank, error) {
	if err := validateTitle(t); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Bank{
		id:        id,
		t:         t,
		questions: []uuid.UUID{},
	}, nil
}

func Restore(id uuid.UUID, t title.Title, questions []uuid.UUID) (*Bank, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateTitle(t); err != nil {
		return nil, err
	}

	if err := validateQuestions(questions); err != nil {
		return nil, err
	}

	return &Bank{
		id:        id,
		t:         t,
		questions: slices.Clone(questions),
	}, nil
}

func (b *Bank) ID() uuid.UUID {
	return b.id
}

func (b *Bank) Title() title.Title {
	return b.t
}

func (b *Bank) Questions() []uuid.UUID {
	return slices.Clone(b.questions)
}

func (b *Bank) CountQuestions() int {
	return len(b.questions)
}

func (b *Bank) HasQuestion(questionID uuid.UUID) bool {
	return slices.Contains(b.questions, questionID)
}

func (b *Bank) Rename(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}

	b.t = t
	return nil
}

func (b *Bank) AddQuestions(questions ...uuid.UUID) error {
	if len(questions) == 0 {
		return nil
	}

	if err := validateQuestionsForAdding(b.questions, questions); err != nil {
		return err
	}

	b.questions = append(b.questions, questions...)
	return nil
}

func (b *Bank) RemoveQuestions(questions ...uuid.UUID) error {
	if len(questions) == 0 {
		return nil
	}

	if err := validateQuestionsForRemoving(b.questions, questions); err != nil {
		return err
	}

	for i := range questions {
		if idx := slices.Index(b.questions, questions[i]); idx != -1 {
			b.questions = slices.Delete(b.questions, idx, idx+1)
		}
	}

	return nil
}

func (b *Bank) ClearQuestions() {
	b.questions = b.questions[:0]
}
