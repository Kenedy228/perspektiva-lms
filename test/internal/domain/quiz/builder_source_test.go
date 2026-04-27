package quiz_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz"
	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type sourceBuilder struct {
	bankID   uuid.UUID
	criteria criteria.Criteria
}

func newSourceBuilder() *sourceBuilder {
	return &sourceBuilder{
		bankID:   uuid.Nil,
		criteria: nil,
	}
}

func (b *sourceBuilder) withBank(id uuid.UUID) *sourceBuilder {
	b.bankID = id
	return b
}

func (b *sourceBuilder) withCriteria(criteria criteria.Criteria) *sourceBuilder {
	b.criteria = criteria
	return b
}

func (b *sourceBuilder) build(t *testing.T, wantErr error) quiz.Source {
	t.Helper()

	s, err := quiz.NewSource(b.bankID, b.criteria)

	assert.ErrorIs(t, err, wantErr)

	return s
}

func (b *sourceBuilder) buildNoTest() quiz.Source {
	s, _ := quiz.NewSource(b.bankID, b.criteria)
	return s
}

func mockSource() quiz.Source {
	criteria := new(mockCriteria)
	s, _ := quiz.NewSource(uuid.New(), criteria)
	return s
}

func mockSourceList(ids ...uuid.UUID) []quiz.Source {
	sources := make([]quiz.Source, 0, len(ids))

	for i := range ids {
		criteria := new(mockCriteria)
		s, _ := quiz.NewSource(ids[i], criteria)
		sources = append(sources, s)
	}

	return sources
}

func mockSourcesWithLength(length int) []quiz.Source {
	sources := make([]quiz.Source, 0, length)

	for range length {
		s := mockSource()
		sources = append(sources, s)
	}

	return sources
}
