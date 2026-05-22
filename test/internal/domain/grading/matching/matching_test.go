//go:build legacy
// +build legacy

package matching_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/matching"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	// Подготавливаем эталонный вопрос (Установление соответствия):
	// Текст: "Сопоставьте столицы с их странами"
	// Пары:
	// 1. Париж -> Франция
	// 2. Берлин -> Германия
	// 3. Токио -> Япония
	q := newQuestionBuilder().
		withTitle("Сопоставьте столицы с их странами").
		withPair("Париж", "Франция").
		withPair("Берлин", "Германия").
		withPair("Токио", "Япония").
		build(t)

	correctPairs := q.Pairs()
	require.Len(t, correctPairs, 3)

	// Извлекаем сгенерированные ID для тестов
	promptParisID := correctPairs[0].PromptID()
	matchFranceID := correctPairs[0].MatchID()

	promptBerlinID := correctPairs[1].PromptID()
	matchGermanyID := correctPairs[1].MatchID()

	promptTokyoID := correctPairs[2].PromptID()
	matchJapanID := correctPairs[2].MatchID()

	fakeID := uuid.New() // Несуществующий ID

	c := matching.New()

	tests := []struct {
		name          string
		studentPairs  map[uuid.UUID]uuid.UUID // map[PromptID]MatchID
		expectedScore float64
	}{
		{
			name: "идеально правильное сопоставление",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchFranceID,
				promptBerlinID: matchGermanyID,
				promptTokyoID:  matchJapanID,
			},
			expectedScore: 1.0,
		},
		{
			name: "перепутаны два ответа",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchGermanyID, // Ошибка здесь
				promptBerlinID: matchFranceID,  // И здесь
				promptTokyoID:  matchJapanID,
			},
			expectedScore: 0.0,
		},
		{
			name: "пропущена одна пара (недобор)",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchFranceID,
				promptBerlinID: matchGermanyID,
			},
			expectedScore: 0.0,
		},
		{
			name: "добавлена лишняя пара (перебор/хак)",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchFranceID,
				promptBerlinID: matchGermanyID,
				promptTokyoID:  matchJapanID,
				fakeID:         fakeID,
			},
			expectedScore: 0.0,
		},
		{
			name: "сопоставление с несуществующим вариантом ответа",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchFranceID,
				promptBerlinID: matchGermanyID,
				promptTokyoID:  fakeID,
			},
			expectedScore: 0.0,
		},
		{
			name: "сопоставление нескольких Prompts с одним и тем же match",
			studentPairs: map[uuid.UUID]uuid.UUID{
				promptParisID:  matchFranceID,
				promptBerlinID: matchFranceID, // Ошибка
				promptTokyoID:  matchJapanID,
			},
			expectedScore: 0.0, // Упадет на проверке Берлина (ожидалась Германия)
		},
		{
			name:          "пустой ответ",
			studentPairs:  map[uuid.UUID]uuid.UUID{},
			expectedScore: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := makeAnswer(tt.studentPairs)
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
			name:     "поддерживает только matching",
			t:        question.TypeMatching,
			supports: true,
		},
		{
			name:     "формат selectable не поддерживается",
			t:        question.TypeSelectable,
			supports: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			c := matching.New()

			//Act
			supports := c.Supports(tt.t)

			//Assert
			assert.Equal(t, tt.supports, supports)
		})
	}
}
