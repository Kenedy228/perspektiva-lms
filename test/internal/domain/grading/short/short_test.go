//go:build legacy
// +build legacy

package short_test

import (
	"testing"

	short2 "gitflic.ru/lms/backend/internal/domain/grading/short"
	"gitflic.ru/lms/backend/internal/domain/question"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupports(t *testing.T) {
	tc := []struct {
		name     string
		t        question.Type
		supports bool
	}{
		{
			name:     "поддерживает только short",
			t:        question.TypeShort,
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
			c := short2.Checker{}

			//Act
			supports := c.Supports(tt.t)

			//Assert
			assert.Equal(t, tt.supports, supports)
		})
	}
}

func TestCheck(t *testing.T) {
	variants := []variant.Variant{
		makeVariant("вариант 1"),
		makeVariant("Вариант 2"),
		makeVariant("ВарИаНТ 3"),
		makeVariant("  вариант 4   "),
	}
	q, err := qshort.New(makeTitle(), variants)
	require.NoError(t, err)

	t.Run("проверка ответа без конфигурации", func(t *testing.T) {
		c := short2.New()

		tc := []struct {
			name    string
			score   float64
			variant string
		}{
			{
				name:    "вариант ответа 2, регистр совпадает",
				score:   1,
				variant: "Вариант 2",
			},
			{
				name:    "вариант ответа 4, пробелы не совпадают",
				score:   0,
				variant: "вариант 4",
			},
			{
				name:    "вариант ответа 3, регистр совпадает",
				score:   1,
				variant: "ВарИаНТ 3",
			},
			{
				name:    "вариант ответа 3, регистр не совпадает",
				score:   0,
				variant: "варИаНТ 3",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				//Arrange
				ans := answer.New(makeAnswerVariant(tt.variant))

				//Act
				s, err := c.Check(q, ans)

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.score, s.Value())
			})
		}
	})

	t.Run("проверка правильного ответа с удалением незначащих пробелов", func(t *testing.T) {
		c := short2.New(short2.TrimSpace())

		tc := []struct {
			name    string
			score   float64
			variant string
		}{
			{
				name:    "вариант ответа 2, в ответе пробелы",
				score:   1,
				variant: "   Вариант 2    ",
			},
			{
				name:    "вариант ответа 4, пробелы не совпадают",
				score:   1,
				variant: "вариант 4",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				//Arrange
				ans := answer.New(makeAnswerVariant(tt.variant))

				//Act
				s, err := c.Check(q, ans)

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.score, s.Value())
			})
		}
	})

	t.Run("проверка правильного ответа без учета регистра", func(t *testing.T) {
		c := short2.New(short2.ToLower())

		tc := []struct {
			name    string
			score   float64
			variant string
		}{
			{
				name:    "вариант ответа 2, регистры не совпадают",
				score:   1,
				variant: "ВАРИАНТ 2",
			},
			{
				name:    "вариант ответа 4, пробелы не совпадают",
				score:   0,
				variant: "   вариант 4           ",
			},
			{
				name:    "вариант ответа 3, регистры не совпадают",
				score:   1,
				variant: "вариант 3",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				//Arrange
				ans := answer.New(makeAnswerVariant(tt.variant))

				//Act
				s, err := c.Check(q, ans)

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.score, s.Value())
			})
		}
	})

	t.Run("проверка правильного ответа без учета регистра и с удалением незначащих пробелов", func(t *testing.T) {
		c := short2.New(short2.ToLower(), short2.TrimSpace())

		tc := []struct {
			name    string
			score   float64
			variant string
		}{
			{
				name:    "вариант ответа 2, регистры не совпадают, дополнительные пробелы",
				score:   1,
				variant: "ВАРИАНТ 2       ",
			},
			{
				name:    "вариант ответа 4, пробелы не совпадают",
				score:   1,
				variant: "   вариант 4           ",
			},
			{
				name:    "вариант ответа 3, регистры не совпадают",
				score:   1,
				variant: "вариант 3",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				//Arrange
				ans := answer.New(makeAnswerVariant(tt.variant))

				//Act
				s, err := c.Check(q, ans)

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.score, s.Value())
			})
		}
	})
}
