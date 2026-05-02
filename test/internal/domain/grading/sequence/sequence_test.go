package sequence_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cseq "gitflic.ru/lms/internal/domain/grading/sequence"
	"gitflic.ru/lms/internal/domain/question"
)

func TestChecker_Check(t *testing.T) {
	q := newQuestionBuilder().
		withTitle("Расположите события в хронологическом порядке").
		withOption("Событие 1 (Сначала)").
		withOption("Событие 2 (Затем)").
		withOption("Событие 3 (В конце)").
		build(t)

	correctOptions := q.Options()
	require.Len(t, correctOptions, 3)

	id1 := correctOptions[0].ID()
	id2 := correctOptions[1].ID()
	id3 := correctOptions[2].ID()
	fakeID := uuid.New() // Несуществующий ID для тестов на ошибку

	c := cseq.New()

	tests := []struct {
		name           string
		studentOptions []uuid.UUID
		expectedScore  float64
	}{
		{
			name:           "идеальный порядок",
			studentOptions: []uuid.UUID{id1, id2, id3},
			expectedScore:  1.0,
		},
		{
			name:           "неправильный порядок (перепутаны последние два)",
			studentOptions: []uuid.UUID{id1, id3, id2},
			expectedScore:  0.0,
		},
		{
			name:           "неправильный порядок (обратный)",
			studentOptions: []uuid.UUID{id3, id2, id1},
			expectedScore:  0.0,
		},
		{
			name:           "ответ содержит меньше элементов",
			studentOptions: []uuid.UUID{id1, id2},
			expectedScore:  0.0,
		},
		{
			name:           "ответ содержит больше элементов",
			studentOptions: []uuid.UUID{id1, id2, id3, fakeID},
			expectedScore:  0.0,
		},
		{
			name:           "ответ содержит несуществующие элементы",
			studentOptions: []uuid.UUID{id1, fakeID, id3},
			expectedScore:  0.0,
		},
		{
			name:           "пустой ответ",
			studentOptions: []uuid.UUID{},
			expectedScore:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := makeAnswer(tt.studentOptions)
			score, err := c.Check(q, ans)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedScore, score.Value())
		})
	}
}

func TestSupports(t *testing.T) {
	tc := []struct {
		name     string
		t        question.Type
		supports bool
	}{
		{
			name:     "поддерживает только sequence",
			t:        question.TypeSequence,
			supports: true,
		},
		{
			name:     "формат matching не поддерживается",
			t:        question.TypeMatching,
			supports: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			c := cseq.New()

			//Act
			supports := c.Supports(tt.t)

			//Assert
			assert.Equal(t, tt.supports, supports)
		})
	}
}
