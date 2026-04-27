package bank

import (
	"slices"
	"time"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Bank struct {
	id        uuid.UUID
	title     string
	questions []uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	deletedAt time.Time
}

func New(title string) (*Bank, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	id, err := uid.New()
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

func (b *Bank) DeletedAt() time.Time {
	return b.deletedAt
}

func (b *Bank) Rename(title string) error {
	if b.IsDeleted() {
		return ErrDeleted
	}

	if err := validateTitle(title); err != nil {
		return err
	}

	b.title = title
	b.updatedAt = time.Now()
	return nil
}

func (b *Bank) AddQuestions(questions ...uuid.UUID) error {
	if b.IsDeleted() {
		return ErrDeleted
	}

	if len(questions) == 0 {
		return nil
	}

	if err := validateQuestionsForAdding(b.questions, questions); err != nil {
		return err
	}

	b.questions = append(b.questions, questions...)
	b.updatedAt = time.Now()
	return nil
}

func (b *Bank) RemoveQuestions(questions ...uuid.UUID) error {
	if b.IsDeleted() {
		return ErrDeleted
	}

	if len(questions) == 0 {
		return nil
	}

	hasDeleted := false
	for i := range questions {
		if idx := slices.Index(b.questions, questions[i]); idx != -1 {
			b.questions = slices.Delete(b.questions, idx, idx+1)
			hasDeleted = true
		}
	}

	if hasDeleted {
		b.updatedAt = time.Now()
	}

	return nil
}

func (b *Bank) ClearQuestions() error {
	if b.IsDeleted() {
		return ErrDeleted
	}

	b.questions = b.questions[:0]
	b.updatedAt = time.Now()
	return nil
}

func (b *Bank) IsDeleted() bool {
	return !b.deletedAt.IsZero()
}

func (b *Bank) Delete() error {
	if b.IsDeleted() {
		return ErrDeleted
	}

	now := time.Now()
	b.deletedAt = now
	b.updatedAt = now
	return nil
}

func (b *Bank) Restore() error {
	if !b.IsDeleted() {
		return ErrNotDeleted
	}

	b.deletedAt = time.Time{}
	b.updatedAt = time.Now()
	return nil
}
