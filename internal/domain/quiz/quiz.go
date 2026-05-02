package quiz

import (
	"slices"

	"gitflic.ru/lms/internal/domain/quiz/limit"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"gitflic.ru/lms/internal/domain/quiz/title"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Quiz struct {
	id       uuid.UUID
	time     limit.Time
	attempts limit.Attempts
	t        title.Title
	sources  []source.Source
}

func New(params Params) (*Quiz, error) {
	if err := validateSources(params.Sources); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Quiz{
		id:       id,
		t:        params.Title,
		sources:  slices.Clone(params.Sources),
		attempts: params.MaxAttempts,
		time:     params.TimeLimit,
	}, nil
}

func (q *Quiz) ID() uuid.UUID {
	return q.id
}

func (q *Quiz) Title() title.Title {
	return q.t
}

func (q *Quiz) Sources() []source.Source {
	return slices.Clone(q.sources)
}

func (q *Quiz) Attempts() limit.Attempts {
	return q.attempts
}

func (q *Quiz) Time() limit.Time {
	return q.time
}

func (q *Quiz) HasInfiniteAttempts() bool {
	return q.Attempts().IsInfinite()
}

func (q *Quiz) HasInfiniteTime() bool {
	return q.time.IsInfinite()
}

func (q *Quiz) ChangeMaxAttempts(maxAttempts limit.Attempts) {
	q.attempts = maxAttempts
}

func (q *Quiz) ChangeTimeLimit(limit limit.Time) {
	q.time = limit
}

func (q *Quiz) Rename(t title.Title) {
	q.t = t
}

func (q *Quiz) AddSource(s source.Source) error {
	if err := validateSourceToAdd(q.sources, s); err != nil {
		return err
	}

	q.sources = append(q.sources, s)
	return nil
}

func (q *Quiz) RemoveSource(s source.Source) error {
	if err := validateSourcesToRemove(q.sources); err != nil {
		return err
	}

	for i := range q.sources {
		if q.sources[i].BankID() == s.BankID() {
			q.sources = slices.Delete(q.sources, i, i+1)
			break
		}
	}

	return nil
}
