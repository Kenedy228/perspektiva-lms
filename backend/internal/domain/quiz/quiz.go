package quiz

import (
	"slices"

	limit2 "gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Quiz struct {
	id               uuid.UUID
	time             limit2.Time
	attempts         limit2.Attempts
	t                title.Title
	shuffleQuestions bool
	sources          []source.Source
}

func New(params Params) (*Quiz, error) {
	if err := validateTitle(params.Title); err != nil {
		return nil, err
	}

	if err := validateSources(params.Sources); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Quiz{
		id:               id,
		t:                params.Title,
		sources:          slices.Clone(params.Sources),
		attempts:         params.MaxAttempts,
		time:             params.TimeLimit,
		shuffleQuestions: params.ShuffleQuestions,
	}, nil
}

func Restore(id uuid.UUID, params Params) (*Quiz, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateTitle(params.Title); err != nil {
		return nil, err
	}

	if err := validateSources(params.Sources); err != nil {
		return nil, err
	}

	return &Quiz{
		id:               id,
		t:                params.Title,
		sources:          slices.Clone(params.Sources),
		attempts:         params.MaxAttempts,
		time:             params.TimeLimit,
		shuffleQuestions: params.ShuffleQuestions,
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

func (q *Quiz) Attempts() limit2.Attempts {
	return q.attempts
}

func (q *Quiz) Time() limit2.Time {
	return q.time
}

func (q *Quiz) ShuffleQuestions() bool {
	return q.shuffleQuestions
}

func (q *Quiz) HasInfiniteAttempts() bool {
	return q.Attempts().IsInfinite()
}

func (q *Quiz) HasInfiniteTime() bool {
	return q.time.IsInfinite()
}

func (q *Quiz) ChangeMaxAttempts(maxAttempts limit2.Attempts) error {
	q.attempts = maxAttempts
	return nil
}

func (q *Quiz) ChangeTimeLimit(limit limit2.Time) error {
	q.time = limit
	return nil
}

func (q *Quiz) ChangeShuffleQuestions(enabled bool) {
	q.shuffleQuestions = enabled
}

func (q *Quiz) Rename(t title.Title) error {
	if err := validateTitle(t); err != nil {
		return err
	}

	q.t = t
	return nil
}

func (q *Quiz) ReplaceSources(sources []source.Source) error {
	if err := validateSources(sources); err != nil {
		return err
	}

	q.sources = slices.Clone(sources)
	return nil
}

func (q *Quiz) ChangeSource(s source.Source) error {
	if s.IsZero() {
		return validateSourceToAdd(q.sources, s)
	}

	if err := validateSourceExists(q.sources, s.BankID()); err != nil {
		return err
	}

	for i := range q.sources {
		if q.sources[i].BankID() == s.BankID() {
			q.sources[i] = s
			break
		}
	}

	return nil
}

func (q *Quiz) AddSource(s source.Source) error {
	if err := validateSourceToAdd(q.sources, s); err != nil {
		return err
	}

	q.sources = append(q.sources, s)
	return nil
}

func (q *Quiz) RemoveSource(s source.Source) error {
	return q.RemoveSourceByBankID(s.BankID())
}

func (q *Quiz) RemoveSourceByBankID(bankID uuid.UUID) error {
	if err := validateSourcesToRemove(q.sources); err != nil {
		return err
	}

	if err := validateSourceExists(q.sources, bankID); err != nil {
		return err
	}

	for i := range q.sources {
		if q.sources[i].BankID() == bankID {
			q.sources = slices.Delete(q.sources, i, i+1)
			break
		}
	}

	return nil
}
