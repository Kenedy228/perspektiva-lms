package typed_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/grading/typed"
	"gitflic.ru/lms/internal/domain/question"
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
			name:     "поддерживает только typed",
			t:        question.TypeTyped,
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
			c := typed.Checker{}

			//Act
			supports := c.Supports(tt.t)

			//Assert
			assert.Equal(t, tt.supports, supports)
		})
	}
}

func TestCheck(t *testing.T) {
	q := newQuestionBuilder().
		withTitle("Столица России - {{moscow}}. Столица Франции - {{paris}}.").
		withBlank("{{moscow}}", "Москва", "г. Москва").
		withBlank("{{paris}}", "Париж").
		build(t)

	t.Run("проверка без нормализаторов (строгая)", func(t *testing.T) {
		c := typed.New()

		tests := []struct {
			name          string
			studentBlanks map[string]string
			expectedScore float64
		}{
			{
				name: "всё абсолютно верно",
				studentBlanks: map[string]string{
					"{{moscow}}": "Москва",
					"{{paris}}":  "Париж",
				},
				expectedScore: 1.0,
			},
			{
				name: "использование альтернативного варианта",
				studentBlanks: map[string]string{
					"{{moscow}}": "г. Москва",
					"{{paris}}":  "Париж",
				},
				expectedScore: 1.0,
			},
			{
				name: "ошибка в регистре (строгая проверка)",
				studentBlanks: map[string]string{
					"{{moscow}}": "москва",
					"{{paris}}":  "Париж",
				},
				expectedScore: 0.0,
			},
			{
				name: "лишние пробелы (строгая проверка)",
				studentBlanks: map[string]string{
					"{{moscow}}": "Москва",
					"{{paris}}":  " Париж ",
				},
				expectedScore: 0.0,
			},
			{
				name: "не заполнен один пропуск",
				studentBlanks: map[string]string{
					"{{moscow}}": "Москва",
				},
				expectedScore: 0.0,
			},
			{
				name: "неправильный ответ в одном из пропусков",
				studentBlanks: map[string]string{
					"{{moscow}}": "Москва",
					"{{paris}}":  "Лондон",
				},
				expectedScore: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ans := makeAnswer(tt.studentBlanks)
				score, err := c.Check(q, ans)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedScore, score.Value())
			})
		}
	})

	t.Run("проверка с обрезкой пробелов (TrimSpace)", func(t *testing.T) {
		c := typed.New(typed.TrimSpace())

		tests := []struct {
			name          string
			studentBlanks map[string]string
			expectedScore float64
		}{
			{
				name: "правильные ответы с пробелами вокруг",
				studentBlanks: map[string]string{
					"{{moscow}}": "   Москва  ",
					"{{paris}}":  "\tПариж\n",
				},
				expectedScore: 1.0,
			},
			{
				name: "ошибка в регистре (TrimSpace не спасет)",
				studentBlanks: map[string]string{
					"{{moscow}}": "москва",
					"{{paris}}":  "Париж",
				},
				expectedScore: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ans := makeAnswer(tt.studentBlanks)
				score, err := c.Check(q, ans)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedScore, score.Value())
			})
		}
	})

	t.Run("проверка с игнорированием регистра (ToLower)", func(t *testing.T) {
		c := typed.New(typed.ToLower())

		tests := []struct {
			name          string
			studentBlanks map[string]string
			expectedScore float64
		}{
			{
				name: "всё верно, разные регистры",
				studentBlanks: map[string]string{
					"{{moscow}}": "МОСКВА",
					"{{paris}}":  "париж",
				},
				expectedScore: 1.0,
			},
			{
				name: "всё верно, альтернативный вариант в другом регистре",
				studentBlanks: map[string]string{
					"{{moscow}}": "Г. мОсКвА",
					"{{paris}}":  "ПаРиЖ",
				},
				expectedScore: 1.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ans := makeAnswer(tt.studentBlanks)
				score, err := c.Check(q, ans)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedScore, score.Value())
			})
		}
	})

	t.Run("проверка с TrimSpace и ToLower одновременно", func(t *testing.T) {
		c := typed.New(typed.TrimSpace(), typed.ToLower())

		tests := []struct {
			name          string
			studentBlanks map[string]string
			expectedScore float64
		}{
			{
				name: "смесь пробелов и разных регистров",
				studentBlanks: map[string]string{
					"{{moscow}}": "   москва   ",
					"{{paris}}":  " ПАРИЖ",
				},
				expectedScore: 1.0,
			},
			{
				name: "неправильный ответ",
				studentBlanks: map[string]string{
					"{{moscow}}": "   питер   ",
					"{{paris}}":  " ПАРИЖ",
				},
				expectedScore: 0.0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ans := makeAnswer(tt.studentBlanks)
				score, err := c.Check(q, ans)

				require.NoError(t, err)
				assert.Equal(t, tt.expectedScore, score.Value())
			})
		}
	})
}
