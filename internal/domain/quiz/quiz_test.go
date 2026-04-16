package quiz

import (
	"errors"
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria/random"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"github.com/google/uuid"
)

func TestNewQuiz(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		err          error
	}{
		{
			name:         "empty title",
			title:        "",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "whitespaces title",
			title:        " ",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptyTitle,
		},
		{
			name:         "empty sources",
			title:        "title",
			sourceCount:  0,
			attemptLimit: 0,
			timeLimit:    0,
			err:          ErrEmptySources,
		},
		{
			name:         "negative attempts",
			title:        "title",
			sourceCount:  1,
			attemptLimit: -1,
			timeLimit:    0,
			err:          ErrNegativeAttempts,
		},
		{
			name:         "negative time",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    -1,
			err:          ErrNegativeTime,
		},
		{
			name:         "valid",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			q, err := NewQuiz(params)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}

			if err == nil {
				if q.Title() != tt.title {
					t.Errorf("expected title %v, got %v", tt.title, q.Title())
				}

				if len(q.Sources()) != tt.sourceCount {
					t.Errorf("expected len of sources %v, got %v", tt.sourceCount, len(q.Sources()))
				}

				if q.AttemptLimit() != tt.attemptLimit {
					t.Errorf("expected attempts %v, got %v", tt.attemptLimit, q.AttemptLimit())
				}

				if q.TimeLimit() != tt.timeLimit {
					t.Errorf("expected time %v, got %v", tt.timeLimit, q.TimeLimit())
				}

				if len(q.Sources()) != 0 {
					if &q.Sources()[0] == &params.Sources[0] {
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
		sourceCount  int
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true,
		},
		{
			name:         "finite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 10,
			timeLimit:    0,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			q, err := NewQuiz(params)
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
		sourceCount  int
		attemptLimit int
		timeLimit    int
		isFinite     bool
	}{
		{
			name:         "infinite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    0,
			isFinite:     true,
		},
		{
			name:         "finite",
			title:        "title",
			sourceCount:  1,
			attemptLimit: 0,
			timeLimit:    10,
			isFinite:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams(tt.title, tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			q, err := NewQuiz(params)

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
		name         string
		title        string
		sourceCount  int
		attemptLimit int
		timeLimit    int
		err          error
	}{
		{
			name:         "empty title",
			title:        "",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          ErrEmptyTitle,
		},
		{
			name:         "whitespaces title",
			title:        " ",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          ErrEmptyTitle,
		},
		{
			name:         "valid title",
			title:        "new title",
			sourceCount:  1,
			attemptLimit: 1,
			timeLimit:    1,
			err:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := createParams("title", tt.sourceCount, tt.attemptLimit, tt.timeLimit)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}

			q, err := NewQuiz(params)
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
	t.Run("successive add", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		toAdd, err := createSource(uuid.Nil)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(toAdd)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		if len(q.Sources()) != 11 {
			t.Errorf("expected len of source %d, got %d", 11, q.Sources())
		}
	})

	t.Run("duplicated source", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		toAdd, err := createSource(uuid.Nil)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(toAdd)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(toAdd)
		if err != ErrDuplicatedSource {
			t.Errorf("expected err %v, got %v", ErrDuplicatedSource, err)
		}
	})

	t.Run("duplicated duplicated bank", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		toAdd, err := createSource(uuid.Nil)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(toAdd)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		dupl, err := createSource(toAdd.BankID())
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(dupl)
		if err != ErrDuplicatedBank {
			t.Errorf("expected err %v, got %v", ErrDuplicatedBank, err)
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("remove unexisting elem", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		toRemove, err := createSource(uuid.Nil)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.RemoveSource(toRemove)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}
	})

	t.Run("remove existing elem", func(t *testing.T) {
		params, err := createParams("title", 10, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		toAdd, err := createSource(uuid.Nil)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.AddSource(toAdd)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.RemoveSource(toAdd)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}
	})

	t.Run("remove last elem", func(t *testing.T) {
		params, err := createParams("title", 1, 1, 1)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		q, err := NewQuiz(params)
		if err != nil {
			t.Errorf("expected no err, got %v", err)
		}

		err = q.RemoveSource(params.Sources[0])
		if err != ErrCannotRemoveLastSource {
			t.Errorf("expected err %v, got %v", ErrCannotRemoveLastSource, err)
		}
	})
}

func TestDelete(t *testing.T) {
	params, err := createParams("title", 1, 1, 1)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	}

	q, err := NewQuiz(params)
	if err != nil {
		t.Errorf("expected no err, got %v", err)
	}

	if q.IsDeleted() != false {
		t.Errorf("expected quiz be not deleted, got true")
	}

	q.Delete()

	if q.IsDeleted() != true {
		t.Errorf("expected quiz be deleted, got false")
	}
}

func createParams(title string, sourceCount, attemptLimit, timeLimit int) (Params, error) {
	sources := make([]source.Source, 0, 5)
	for range sourceCount {
		s, err := createSource(uuid.Nil)
		if err != nil {
			return Params{}, err
		}

		sources = append(sources, s)
	}

	return Params{
		Title:        title,
		Sources:      sources,
		AttemptLimit: attemptLimit,
		TimeLimit:    timeLimit,
	}, nil
}

func createSource(bankID uuid.UUID) (source.Source, error) {
	cParams := random.Params{
		Count: 10,
	}

	r, err := random.NewRandomCriteria(cParams)
	if err != nil {
		return source.Source{}, err
	}

	var params source.Params

	if bankID != uuid.Nil {
		params = source.Params{
			BankID:   bankID,
			Criteria: r,
		}
	} else {
		params = source.Params{
			BankID:   uuid.New(),
			Criteria: r,
		}
	}

	s, err := source.NewSource(params)
	if err != nil {
		return source.Source{}, err
	}

	return s, nil
}
