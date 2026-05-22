//go:build legacy
// +build legacy

package quiz_test

import (
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockCriteria struct {
	mock.Mock
}

func (m *mockCriteria) QuestionCount() int  { return 10 }
func (m *mockCriteria) Type() criteria.Type { return criteria.TypeRandom }

func mockSource() source.Source {
	s, _ := source.NewSource(uuid.New(), new(mockCriteria))
	return s
}

func mockSourceList(ids ...uuid.UUID) []source.Source {
	sources := make([]source.Source, 0, len(ids))
	for i := range ids {
		s, _ := source.NewSource(ids[i], new(mockCriteria))
		sources = append(sources, s)
	}
	return sources
}

func mockSourcesWithLength(length int) []source.Source {
	sources := make([]source.Source, 0, length)
	for i := 0; i < length; i++ {
		sources = append(sources, mockSource())
	}
	return sources
}
