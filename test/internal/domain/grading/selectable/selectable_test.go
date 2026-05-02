package selectable_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cselectable "gitflic.ru/lms/internal/domain/grading/selectable"
	"gitflic.ru/lms/internal/domain/question"
)

func TestChecker_Check(t *testing.T) {
	q := newQuestionBuilder().
		withTitle("Какие из этих городов находятся в Европе?").
		withOption("Париж", true).
		withOption("Токио", false).
		withOption("Берлин", true).
		withOption("Нью-Йорк", false).
		build(t)

	correctOptions := q.Options()
	require.Len(t, correctOptions, 4)

	// Извлекаем ID для тестов
	parisID := correctOptions[0].ID()   // правильный
	tokyoID := correctOptions[1].ID()   // неправильный
	berlinID := correctOptions[2].ID()  // правильный
	newYorkID := correctOptions[3].ID() // неправильный
	fakeID := uuid.New()                // несуществующий в вопросе ID

	c := cselectable.New()

	tests := []struct {
		name           string
		studentOptions []uuid.UUID
		expectedScore  float64
	}{
		{
			name:           "выбраны все правильные опции (идеально)",
			studentOptions: []uuid.UUID{parisID, berlinID},
			expectedScore:  1.0,
		},
		{
			name:           "выбрана только часть правильных опций (недобор)",
			studentOptions: []uuid.UUID{parisID},
			expectedScore:  0.0,
		},
		{
			name:           "выбраны все правильные и один неправильный (перебор/штраф)",
			studentOptions: []uuid.UUID{parisID, berlinID, tokyoID},
			expectedScore:  0.0,
		},
		{
			name:           "выбраны только неправильные",
			studentOptions: []uuid.UUID{tokyoID, newYorkID},
			expectedScore:  0.0,
		},
		{
			name:           "пустой ответ",
			studentOptions: []uuid.UUID{},
			expectedScore:  0.0,
		},
		{
			name:           "ответ содержит несуществующий ID (игнорируется твоим кодом, но итоговая проверка падает по счетчику)",
			studentOptions: []uuid.UUID{parisID, berlinID, fakeID},
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
			name:     "поддерживает только selectable",
			t:        question.TypeSelectable,
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
			c := cselectable.New()

			//Act
			supports := c.Supports(tt.t)

			//Assert
			assert.Equal(t, tt.supports, supports)
		})
	}
}
