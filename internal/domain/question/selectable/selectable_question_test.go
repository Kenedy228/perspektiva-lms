package selectable

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// === ТЕСТЫ ДЛЯ OPTION ===

func TestNewOption(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		isCorrect bool
		wantErr   error
	}{
		{
			name:      "valid option",
			text:      "Правильный ответ",
			isCorrect: true,
			wantErr:   nil,
		},
		{
			name:      "error empty text",
			text:      "   ",
			isCorrect: false,
			wantErr:   ErrEmptyOptionText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := NewOption(tt.text, tt.isCorrect)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				assert.NotEqual(t, uuid.Nil, opt.ID())
				assert.Equal(t, tt.text, opt.Text())
				assert.Equal(t, tt.isCorrect, opt.IsCorrect())
			}
		})
	}
}

// === ТЕСТЫ ДЛЯ SELECTABLE QUESTION ===

func TestNewSelectableQuestion(t *testing.T) {
	// Генерируем мапу с 21 элементом для проверки лимита (maxOptions = 20)
	tooManyOptions := make(map[string]bool, maxOptions+1)
	for i := 0; i <= maxOptions; i++ {
		tooManyOptions[fmt.Sprintf("option %d", i)] = true
	}

	tests := []struct {
		name    string
		params  *Params
		wantErr error
	}{
		{
			name: "success valid question",
			params: &Params{
				Text: "Выберите цвета светофора",
				Options: map[string]bool{
					"Красный": true,
					"Зеленый": true,
					"Фиолетовый": false,
				},
			},
			wantErr: nil,
		},
		{
			name: "error empty options",
			params: &Params{
				Text: "Вопрос без вариантов",
				Options: map[string]bool{},
			},
			wantErr: ErrEmptyOptions,
		},
		{
			name: "error not enough options",
			params: &Params{
				Text: "Один вариант ответа",
				Options: map[string]bool{
					"Одинокий ответ": true,
				},
			},
			wantErr: ErrNotEnoughOptions, // minOptions = 2
		},
		{
			name: "error too many options",
			params: &Params{
				Text: "Слишком много вариантов",
				Options: tooManyOptions,
			},
			wantErr: ErrTooManyOptions,
		},
		{
			name: "error no correct options",
			params: &Params{
				Text: "Где правильный?",
				Options: map[string]bool{
					"Неправильно": false,
					"Тоже нет": false,
				},
			},
			wantErr: ErrNoCorrectOption,
		},
		{
			name: "error empty option text",
			params: &Params{
				Text: "Пустой текст в мапе",
				Options: map[string]bool{
					"Нормальный": true,
					"   ": false, // вызовет ошибку валидации в NewOption
				},
			},
			wantErr: ErrEmptyOptionText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewSelectableQuestion(tt.params)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.NotNil(t, q)
				selectableQ, ok := q.(*SelectableQuestion)
				require.True(t, ok)

				assert.Equal(t, question.TypeSelectable, selectableQ.Type())
				assert.Len(t, selectableQ.Options(), len(tt.params.Options))
			}
		})
	}
}

func TestSelectableQuestion_UpdateOptions(t *testing.T) {
	params := &Params{
		Text: "Изначальный текст",
		Options: map[string]bool{
			"Вариант 1": true,
			"Вариант 2": false,
		},
	}
	q, err := NewSelectableQuestion(params)
	require.NoError(t, err)

	selectableQ, ok := q.(*SelectableQuestion)
	require.True(t, ok)

	t.Run("success update options", func(t *testing.T) {
		newOptions := map[string]bool{
			"Новый 1": true,
			"Новый 2": true,
			"Новый 3": false,
		}

		err := selectableQ.UpdateOptions(newOptions)
		
		assert.NoError(t, err)
		assert.Len(t, selectableQ.Options(), 3)
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldOptions := selectableQ.Options()

		// Пытаемся обновить невалидными данными (отсутствует правильный ответ)
		invalidOptions := map[string]bool{
			"Ошибка 1": false,
			"Ошибка 2": false,
		}
		err := selectableQ.UpdateOptions(invalidOptions)

		assert.ErrorIs(t, err, ErrNoCorrectOption)
		
		// Убеждаемся, что старые варианты не затерлись
		assert.Equal(t, oldOptions, selectableQ.Options())
	})
}

func TestSelectableQuestion_Encapsulation(t *testing.T) {
	params := &Params{
		Text: "Текст",
		Options: map[string]bool{
			"Опция 1": true,
			"Опция 2": false,
		},
	}

	q, err := NewSelectableQuestion(params)
	require.NoError(t, err)
	selectableQ := q.(*SelectableQuestion)

	t.Run("getter clone", func(t *testing.T) {
		gotOptions := selectableQ.Options()
		require.Len(t, gotOptions, 2)

		// Пытаемся затереть первый элемент в полученном срезе
		gotOptions[0] = Option{}

		// Проверяем, что внутреннее состояние не изменилось (оригинал не пострадал)
		assert.NotEqual(t, uuid.Nil, selectableQ.Options()[0].ID())
	})
}
