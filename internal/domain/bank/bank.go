package bank

import (
	"slices"
	"time"

	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Bank struct {
	id        uuid.UUID
	title     string
	questions []uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func New(title string) (*Bank, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	id, err := utils.GenerateID()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Bank{
		id:        id,
		title:     title,
		questions: []uuid.UUID{},
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (b *Bank) ID() uuid.UUID {
	return b.id
}

func (b *Bank) Title() string {
	return b.title
}

func (b *Bank) Questions() []uuid.UUID {
	return slices.Clone(b.questions)
}

func (b *Bank) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Bank) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *Bank) Rename(title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}

	b.title = title
	b.updatedAt = time.Now()

	return nil
}

func (b *Bank) AddQuestions(questions ...uuid.UUID) (int, error) {
	if err := validateQuestionsForAdding(b.questions, questions); err != nil {
		return 0, err
	}

	b.questions = append(b.questions, questions...)
	b.updatedAt = time.Now()

	return len(questions), nil
}

func (b *Bank) RemoveQuestions(questions ...uuid.UUID) int {
	total := 0
	for i := range questions {
		index := slices.Index(b.questions, questions[i])
		if index != -1 {
			b.questions = slices.Delete(b.questions, index, index+1)
			total++
		}
	}

	b.updatedAt = time.Now()
	return total
}

func (b *Bank) ClearQuestions() int {
	total := len(b.questions)
	b.questions = b.questions[:0]
	b.updatedAt = time.Now()

	return total
}
