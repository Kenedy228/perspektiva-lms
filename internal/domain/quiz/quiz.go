package quiz

import (
	"slices"
	"time"

	"gitflic.ru/lms/internal/domain/shared/limit"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Quiz struct {
	id          uuid.UUID
	title       string
	sources     []Source
	maxAttempts int
	timeLimit   limit.Limit
	createdAt   time.Time
	updatedAt   time.Time
}

func New(params Params) (*Quiz, error) {
	if err := validateTitle(params.Title); err != nil {
		return nil, err
	}

	if err := validateSources(params.Sources); err != nil {
		return nil, err
	}

	if err := validateMaxAttempts(params.MaxAttempts); err != nil {
		return nil, err
	}

	sCopy := slices.Clone(params.Sources)

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Quiz{
		id:          id,
		title:       params.Title,
		sources:     sCopy,
		maxAttempts: params.MaxAttempts,
		timeLimit:   params.TimeLimit,
		createdAt:   now,
		updatedAt:   now,
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

func (q *Quiz) MaxAttempts() int {
	return q.maxAttempts
}

func (q *Quiz) TimeLimit() limit.Limit {
	return q.timeLimit
}

func (q *Quiz) CreatedAt() time.Time {
	return q.createdAt
}

func (q *Quiz) UpdatedAt() time.Time {
	return q.updatedAt
}

func (q *Quiz) HasInfiniteAttempts() bool {
	return q.maxAttempts == 0
}

func (q *Quiz) HasInfiniteTime() bool {
	return q.timeLimit.IsInfinite()
}

func (q *Quiz) ChangeMaxAttempts(maxAttempts int) error {
	if err := validateMaxAttempts(maxAttempts); err != nil {
		return err
	}

	q.maxAttempts = maxAttempts
	q.updatedAt = time.Now()
	return nil
}

func (q *Quiz) ChangeTimeLimit(limit limit.Limit) {
	q.timeLimit = limit
	q.updatedAt = time.Now()
}

func (q *Quiz) Rename(title string) error {
	if err := validateTitle(title); err != nil {
		return err
	}

	q.title = title
	q.updatedAt = time.Now()
	return nil
}

func (q *Quiz) AddSource(s Source) error {
	if err := validateSourceToAdd(q.sources, s); err != nil {
		return err
	}

	q.sources = append(q.sources, s)
	q.updatedAt = time.Now()
	return nil
}

func (q *Quiz) RemoveSource(s Source) error {
	if len(q.sources) == 1 {
		return ErrCannotRemoveLastSource
	}

	for i := range q.sources {
		if q.sources[i].BankID() == s.BankID() {
			q.sources = slices.Delete(q.sources, i, i+1)
			q.updatedAt = time.Now()
			break
		}
	}

	return nil
}
