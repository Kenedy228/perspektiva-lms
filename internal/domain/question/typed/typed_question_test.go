package typed

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === ТЕСТЫ ДЛЯ BLANK ===

func TestNewBlank(t *testing.T) {
	tests := []struct {
		name    string
		mark    string
		answers []string
		wantErr error
	}{
		{
			name:    "valid blank",
			mark:    "mark1",
			answers: []string{"answer1", "answer2"},
			wantErr: nil,
		},
		{
			name:    "error empty mark",
			mark:    "   ",
			answers: []string{"answer"},
			wantErr: ErrEmptyMark,
		},
		{
			name:    "error no answers",
			mark:    "mark1",
			answers: []string{},
			wantErr: ErrNoBlankAnswers,
		},
		{
			name:    "error empty answer in slice",
			mark:    "mark1",
			answers: []string{"answer1", "  "},
			wantErr: ErrEmptyBlankAnswer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBlank(tt.mark, tt.answers)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				assert.NotEqual(t, uuid.Nil, b.ID())
				assert.Equal(t, tt.mark, b.Mark())
				assert.Equal(t, tt.answers, b.Answers())
			}
		})
	}
}

func TestBlank_Encapsulation(t *testing.T) {
	originalAnswers := []string{"ans1", "ans2"}
	b, err := NewBlank("mark", originalAnswers)
	require.NoError(t, err)

	// Проверяем, что конструктор скопировал слайс
	originalAnswers[0] = "mutated"
	assert.Equal(t, "ans1", b.Answers()[0], "NewBlank should clone the answers slice")

	// Проверяем, что геттер возвращает копию
	returnedAnswers := b.Answers()
	returnedAnswers[0] = "mutated again"
	assert.Equal(t, "ans1", b.Answers()[0], "Answers() getter should return a clone")
}

// === ТЕСТЫ ДЛЯ TYPED QUESTION ===

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		params  *Params
		wantErr error
	}{
		{
			name: "success valid question",
			params: &Params{
				Text:              "Столица России — [moscow], а Франции — [paris].",
				PlaceholdersCount: 2,
				Blanks: map[string][]string{
					"moscow": {"Москва", "г. Москва"},
					"paris":  {"Париж"},
				},
			},
			wantErr: nil,
		},
		{
			name: "error no placeholders count",
			params: &Params{
				Text:              "Текст без пропусков",
				PlaceholdersCount: 0,
				Blanks:            map[string][]string{},
			},
			wantErr: ErrNoPlaceholders,
		},
		{
			name: "error too many placeholders",
			params: &Params{
				Text:              "...", // Текст не важен, валидация упадет раньше
				PlaceholdersCount: 21,
				Blanks:            map[string][]string{},
			},
			wantErr: ErrTooManyPlaceholders,
		},
		{
			name: "error placeholders count mismatch",
			params: &Params{
				Text:              "Столица — [moscow].",
				PlaceholdersCount: 2, // Указали 2, а бланк передаем 1
				Blanks: map[string][]string{
					"moscow": {"Москва"},
				},
			},
			wantErr: ErrPlaceholderCountMismatch,
		},
		{
			name: "error placeholder missing in text",
			params: &Params{
				Text:              "Столица — [moscow].",
				PlaceholdersCount: 1,
				Blanks: map[string][]string{
					"paris": {"Париж"}, // В тексте нет [paris]
				},
			},
			wantErr: ErrPlaceholderMissing,
		},
		{
			name: "error duplicate mark in text",
			params: &Params{
				Text:              "Столица — [moscow]. И еще раз [moscow].",
				PlaceholdersCount: 1,
				Blanks: map[string][]string{
					"moscow": {"Москва"},
				},
			},
			wantErr: ErrMarkDuplicate,
		},
		{
			name: "error invalid blank (propagated from NewBlank)",
			params: &Params{
				Text:              "Столица — [moscow].",
				PlaceholdersCount: 1,
				Blanks: map[string][]string{
					"moscow": {}, // Пустой список ответов
				},
			},
			wantErr: ErrNoBlankAnswers, // Ошибка всплывает из NewBlank
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.NotNil(t, q)
				typedQ, ok := q.(*TypedQuestion)
				require.True(t, ok)

				assert.Len(t, typedQ.Blanks(), tt.params.PlaceholdersCount)
			}
		})
	}
}

func TestTypedQuestion_ReplaceContent(t *testing.T) {
	// Подготовка валидного вопроса
	params := &Params{
		Text:              "Это [first] тест.",
		PlaceholdersCount: 1,
		Blanks: map[string][]string{
			"first": {"первый"},
		},
	}
	q, err := New(params)
	require.NoError(t, err)

	typedQ := q.(*TypedQuestion)

	t.Run("success replace", func(t *testing.T) {
		err := typedQ.ReplaceContent(
			"Это [new] контент",
			1,
			map[string][]string{"new": {"новый"}},
		)

		assert.NoError(t, err)
		assert.Equal(t, "Это [new] контент", typedQ.Text())
		require.Len(t, typedQ.Blanks(), 1)
		assert.Equal(t, "new", typedQ.Blanks()[0].Mark())
	})

	t.Run("error invalid replace leaves state untouched", func(t *testing.T) {
		oldText := typedQ.Text()
		oldBlanks := typedQ.Blanks()

		// Пытаемся заменить на невалидные данные (отсутствует плейсхолдер)
		err := typedQ.ReplaceContent(
			"Текст без марков",
			1,
			map[string][]string{"missing": {"value"}},
		)

		assert.ErrorIs(t, err, ErrPlaceholderMissing)

		// Состояние вопроса не должно было измениться
		assert.Equal(t, oldText, typedQ.Text())
		assert.Equal(t, oldBlanks, typedQ.Blanks())
	})
}

func TestTypedQuestion_Type(t *testing.T) {
	params := &Params{
		Text:              "Текст [m]",
		PlaceholdersCount: 1,
		Blanks:            map[string][]string{"m": {"1"}},
	}
	q, err := New(params)
	require.NoError(t, err)

	typedQ := q.(*TypedQuestion)
	assert.Equal(t, question.TypeTyped, typedQ.Type())
}

func TestTypedQuestion_BlanksEncapsulation(t *testing.T) {
	params := &Params{
		Text:              "Текст [m]",
		PlaceholdersCount: 1,
		Blanks:            map[string][]string{"m": {"1"}},
	}
	q, err := New(params)
	require.NoError(t, err)

	typedQ := q.(*TypedQuestion)

	blanks := typedQ.Blanks()
	require.Len(t, blanks, 1)

	// Изменяем полученный массив
	blanks[0] = Blank{}

	// Проверяем, что оригинал внутри TypedQuestion не пострадал
	assert.NotEqual(t, blanks[0].Mark(), typedQ.Blanks()[0].Mark())
}
