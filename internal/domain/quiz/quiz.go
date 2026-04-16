package quiz

import (
	"slices"
	"strings"
	"time"

	"gitflic.ru/lms/internal/domain/quiz/source"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Quiz struct {
	id           uuid.UUID
	title        string
	sources      []source.Source
	attemptLimit int
	timeLimit    int
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    time.Time
}

func NewQuiz(params Params) (*Quiz, error) {
	if err := validateTitle(params.Title); err != nil {
		return nil, err
	}
	if err := validateSources(params.Sources); err != nil {
		return nil, err
	}
	if err := validateAttemptLimit(params.AttemptLimit); err != nil {
		return nil, err
	}
	if err := validateTimeLimit(params.TimeLimit); err != nil {
		return nil, err
	}

	sCopy := slices.Clone(params.Sources)

	id, err := utils.GenerateID()

	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Quiz{
		id:           id,
		title:        params.Title,
		sources:      sCopy,
		attemptLimit: params.AttemptLimit,
		timeLimit:    params.TimeLimit,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (q *Quiz) ID() uuid.UUID {
	return q.id
}

func (q *Quiz) Title() string {
	return q.title
}

func (q *Quiz) Sources() []Source {
	return slices.Clone(q.sources)
}

func (q *Quiz) AttemptLimit() int {
	return q.attemptLimit
}

func (q *Quiz) TimeLimit() int {
	return q.timeLimit
}

func (q *Quiz) CreatedAt() time.Time {
	return q.createdAt
}

func (q *Quiz) UpdatedAt() time.Time {
	return q.updatedAt
}

func (q *Quiz) DeletedAt() time.Time {
	return q.deletedAt
}

func (q *Quiz) IsInfiniteAttempts() bool {
	return q.attemptLimit == 0
}

func (q *Quiz) IsInfiniteTime() bool {
	return q.timeLimit == 0
}

func (q *Quiz) Rename(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	q.title = title
	q.updatedAt = time.Now()
	return nil
}

func (q *Quiz) AddSource(source source.Source) error {
		

	q.sources = append(q.sources, source)
	q.updatedAt = time.Now()
	return nil
}

func (q *Quiz) RemoveSource(id uuid.UUID) error {
	idx := slices.Index(q.sourceIDs, id)

	if idx != -1 {
		if len(q.sourceIDs) == 1 {
			return ErrCannotRemoveLastSource
		}

		q.sourceIDs = slices.Delete(q.sourceIDs, idx, idx+1)
		q.updatedAt = time.Now()
	}

	return nil
}

func (q *Quiz) IsDeleted() bool {
	return q.deletedAt.IsZero()
}
