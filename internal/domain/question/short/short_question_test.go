package short

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// Генерируем массив из 21 элемента для проверки лимита
	tooManyAnswers := make([]string, maxAnswers+1)
	for i := range tooManyAnswers {
		tooManyAnswers[i] = fmt.Sprintf("answer %d", i)
	}

	tests := []struct {
		name    string
		params  *Params
		wantErr error
	}{
		{
			name: "success valid short question",
			params: &Params{
				Text:           "Назовите столицу Франции",
				Answers:        []string{"Париж", "париж"},
				AllowDuplicate: false,
			},
			wantErr: nil,
		},
		{
			name: "success with duplicates allowed",
			params: &Params{
				Text:           "Тест дубликатов",
				Answers:        []string{"Да", "Да"},
				AllowDuplicate: true,
			},
			wantErr: nil,
		},
		{
			name: "error duplicate answers not allowed",
			params: &Params{
				Text:           "Тест дубликатов",
				Answers:        []string{"Нет", "Нет"},
				AllowDuplicate: false,
			},
			wantErr: ErrDuplicateAnswer,
		},
		{
			name: "error no answers",
			params: &Params{
				Text:           "Вопрос без ответа",
				Answers:        []string{},
				AllowDuplicate: false,
			},
			wantErr: ErrNoAnswers,
		},
		{
			name: "error empty answer string",
			params: &Params{
				Text:           "Пустой ответ",
				Answers:        []string{"Норм ответ", "   "},
				AllowDuplicate: false,
			},
			wantErr: ErrEmptyAnswer,
		},
		{
			name: "error too many answers",
			params: &Params{
				Text:           "Слишком много вариантов",
				Answers:        tooManyAnswers,
				AllowDuplicate: true,
			},
			wantErr: ErrTooManyAnswers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.NotNil(t, q)
				shortQ, ok := q.(*ShortQuestion)
				require.True(t, ok)

				assert.Equal(t, tt.params.Answers, shortQ.Answers())
				assert.Equal(t, question.TypeShort, shortQ.Type())
			}
		})
	}
}

func TestShortQuestion_UpdateAnswers(t *testing.T) {
	params := &Params{
		Text:           "Базовый текст",
		Answers:        []string{"Ответ 1"},
		AllowDuplicate: false,
	}
	q, err := New(params)
	require.NoError(t, err)

	shortQ, ok := q.(*ShortQuestion)
	require.True(t, ok)

	t.Run("success update answers", func(t *testing.T) {
		newAnswers := []string{"Новый 1", "Новый 2"}
		
		err := shortQ.UpdateAnswers(newAnswers, false)
		
		assert.NoError(t, err)
		assert.Equal(t, newAnswers, shortQ.Answers())
		// q.Touch() вызывается внутри, но т.к. base.Base за рамками этого пакета,
		// мы полагаемся на то, что метод отработал без паники.
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldAnswers := shortQ.Answers()

		// Пытаемся обновить невалидными данными (пустой ответ)
		invalidAnswers := []string{"Норм", ""}
		err := shortQ.UpdateAnswers(invalidAnswers, false)

		assert.ErrorIs(t, err, ErrEmptyAnswer)
		// Убеждаемся, что старые ответы не затерлись
		assert.Equal(t, oldAnswers, shortQ.Answers())
	})
}

func TestShortQuestion_Encapsulation(t *testing.T) {
	t.Run("input params clone", func(t *testing.T) {
		inputAnswers := []string{"Оригинал"}
		params := &Params{
			Text:           "Текст",
			Answers:        inputAnswers,
			AllowDuplicate: false,
		}

		q, err := New(params)
		require.NoError(t, err)
		shortQ := q.(*ShortQuestion)

		// Меняем исходный массив, переданный в Params
		inputAnswers[0] = "Мутация"

		// Проверяем, что внутри структуры остался оригинал
		assert.Equal(t, "Оригинал", shortQ.Answers()[0])
	})

	t.Run("getter clone", func(t *testing.T) {
		params := &Params{
			Text:           "Текст",
			Answers:        []string{"Оригинал"},
			AllowDuplicate: false,
		}

		q, err := New(params)
		require.NoError(t, err)
		shortQ := q.(*ShortQuestion)

		// Получаем ответы и пытаемся их изменить
		gotAnswers := shortQ.Answers()
		gotAnswers[0] = "Взлом"

		// Проверяем, что внутреннее состояние не изменилось
		assert.Equal(t, "Оригинал", shortQ.Answers()[0])
	})

	t.Run("update clone", func(t *testing.T) {
		params := &Params{
			Text:           "Текст",
			Answers:        []string{"Оригинал"},
			AllowDuplicate: false,
		}

		q, err := New(params)
		require.NoError(t, err)
		shortQ := q.(*ShortQuestion)

		newAnswers := []string{"Обновление"}
		err = shortQ.UpdateAnswers(newAnswers, false)
		require.NoError(t, err)

		// Изменяем массив после передачи его в UpdateAnswers
		newAnswers[0] = "Взлом 2"

		// Проверяем, что слайс был скопирован при апдейте
		assert.Equal(t, "Обновление", shortQ.Answers()[0])
	})
}
