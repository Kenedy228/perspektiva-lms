package quiz

import (
	"errors"
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestNewQuiz(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sources      []uuid.UUID
		attemptLimit int
		timeLimit    int
		err          error
	}{
		{
			name:         "empty title",
			title:        "",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "whitespaces title",
			title:        " ",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "empty sources",
			title:        "title",
			sources:      nil,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptySources,
		},
		{
			name:         "empty sources",
			title:        "title",
			sources:      []uuid.UUID{uuid.Nil},
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrNilSource,
		},
		{
			name:         "negative attempts",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: -1,
			timeLimit:    0,
			err:          ErrNegativeAttempts,
		},
		{
			name:         "negative time",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    -1,
			err:          ErrNegativeTime,
		},
		{
			name:         "valid",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    0,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz(tt.title, tt.sources, tt.attemptLimit, tt.timeLimit)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if q.Title() != tt.title {
					t.Errorf("expected title %v, got %v", tt.title, q.Title())
				}

				if len(q.SourceIDs()) != len(tt.sources) {
					t.Errorf("expected len of sources %v, got %v", len(tt.sources), len(q.SourceIDs()))
				}

				if q.AttemptLimit() != tt.attemptLimit {
					t.Errorf("expected attempts %v, got %v", tt.attemptLimit, q.AttemptLimit())
				}

				if q.TimeLimit() != tt.timeLimit {
					t.Errorf("expected time %v, got %v", tt.timeLimit, q.TimeLimit())
				}

				if len(q.SourceIDs()) != 0 {
					if &q.SourceIDs()[0] == &tt.sources[0] {
						t.Errorf("expected different slices, got equal")
					}
				}
			}
		})
	}
}

func TestIsFiniteAttempts(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sources      []uuid.UUID
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true,
		},
		{
			name:         "finite",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 10,
			timeLimit:    0,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz(tt.title, tt.sources, tt.attemptLimit, tt.timeLimit)

			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			if q.IsInfiniteAttempts() != tt.isFinite {
				t.Errorf("expected IsFinite %v, got %v", tt.isFinite, q.IsInfiniteAttempts())
			}
		})
	}
}

func TestIsFinitTime(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sources      []uuid.UUID
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true,
		},
		{
			name:         "finite",
			title:        "title",
			sources:      []uuid.UUID{uuid.New()},
			attemptLimit: 0,
			timeLimit:    10,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz(tt.title, tt.sources, tt.attemptLimit, tt.timeLimit)

			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			if q.IsInfiniteTime() != tt.isFinite {
				t.Errorf("expected IsFinit %v, got %v", tt.isFinite, q.IsInfiniteTime())
			}
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name  string
		title string
		err   error
	}{
		{
			name:  "empty title",
			title: "",
			err:   ErrEmptyTitle,
		},
		{
			name:  "whitespaces title",
			title: " ",
			err:   ErrEmptyTitle,
		},
		{
			name:  "valid title",
			title: "new title",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz("title", []uuid.UUID{uuid.New()}, 0, 0)

			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			err = q.Rename(tt.title)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if q.Title() != tt.title {
					t.Errorf("expected rename, got %v", q.Title())
				}
			}
		})
	}
}

func TestAddSource(t *testing.T) {
	tests := []struct {
		name   string
		source uuid.UUID
		err    error
	}{
		{
			name:   "nil source",
			source: uuid.Nil,
			err:    ErrNilSource,
		},
		{
			name:   "valid source",
			source: uuid.New(),
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz("title", []uuid.UUID{uuid.New()}, 0, 0)

			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			length := len(q.SourceIDs())
			err = q.AddSource(tt.source)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if len(q.SourceIDs()) != length+1 {
					t.Errorf("expected length %d, got %d", length+1, len(q.SourceIDs()))
				}
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name     string
		source   uuid.UUID
		toRemove bool
	}{
		{
			name:     "nil source",
			source:   uuid.Nil,
			toRemove: false,
		},
		{
			name:     "valid source",
			source:   uuid.New(),
			toRemove: false,
		},
		{
			name:     "valid source",
			source:   uuid.New(),
			toRemove: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewQuiz("title", []uuid.UUID{uuid.New()}, 0, 0)

			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			if tt.toRemove {
				err := q.AddSource(tt.source)

				if err != nil {
					t.Errorf("expected err nil, got %v", err)
				}
			}

			err = q.RemoveSource(tt.source)
			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			if slices.Contains(q.SourceIDs(), tt.source) {
				t.Errorf("expected to remove element")
			}
		})
	}
}

func TestRemoveWithLastSource(t *testing.T) {
	toRemove := uuid.New()
	q, err := NewQuiz("title", []uuid.UUID{toRemove}, 0, 0)
	if err != nil {
		t.Errorf("expected err nil, got %v", err)
	}

	err = q.RemoveSource(toRemove)
	if !errors.Is(err, ErrCannotRemoveLastSource) {
		t.Errorf("expected err %v, got %v", ErrCannotRemoveLastSource, err)
	}
}
