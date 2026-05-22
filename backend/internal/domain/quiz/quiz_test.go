package quiz

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
	"github.com/google/uuid"
)

func TestNewRequiresTitleAndSources(t *testing.T) {
	_, err := New(Params{Sources: []source.Source{mustSource(t)}})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid quiz error for missing title, got %v", err)
	}

	_, err = New(Params{Title: mustTitle(t)})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid quiz error for missing sources, got %v", err)
	}
}

func TestNewPreservesShufflePolicy(t *testing.T) {
	q, err := New(Params{
		Title:            mustTitle(t),
		ShuffleQuestions: true,
		Sources:          []source.Source{mustSource(t)},
	})
	if err != nil {
		t.Fatalf("create quiz: %v", err)
	}

	if !q.ShuffleQuestions() {
		t.Fatal("expected enabled shuffle questions policy")
	}
}

func TestRestoreRequiresID(t *testing.T) {
	_, err := Restore(uuid.Nil, Params{
		Title:   mustTitle(t),
		Sources: []source.Source{mustSource(t)},
	})
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid quiz error for missing id, got %v", err)
	}
}

func TestReplaceSourcesValidatesWholeSet(t *testing.T) {
	q := mustQuiz(t, mustSource(t))

	err := q.ReplaceSources(nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid quiz error for empty replacement, got %v", err)
	}
}

func TestRemoveSourceRequiresExistingSource(t *testing.T) {
	first := mustSource(t)
	second := mustSource(t)
	q := mustQuiz(t, first, second)

	err := q.RemoveSourceByBankID(uuid.New())
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected invalid quiz error for missing source, got %v", err)
	}
}

func TestTimeLimitErrorsWrapLimitInvalid(t *testing.T) {
	_, err := limit.NewTime(-1)
	if !errors.Is(err, limit.ErrInvalid) {
		t.Fatalf("expected wrapped limit invalid error, got %v", err)
	}
}

func mustQuiz(t *testing.T, sources ...source.Source) *Quiz {
	t.Helper()

	q, err := New(Params{Title: mustTitle(t), Sources: sources})
	if err != nil {
		t.Fatalf("create quiz: %v", err)
	}
	return q
}

func mustTitle(t *testing.T) title.Title {
	t.Helper()

	tt, err := title.New("Final test")
	if err != nil {
		t.Fatalf("create title: %v", err)
	}
	return tt
}

func mustSource(t *testing.T) source.Source {
	t.Helper()

	c, err := criteria.NewRandom(3)
	if err != nil {
		t.Fatalf("create criteria: %v", err)
	}
	s, err := source.NewSource(uuid.New(), c)
	if err != nil {
		t.Fatalf("create source: %v", err)
	}
	return s
}
