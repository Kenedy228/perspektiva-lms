package short

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tooManyAnswers := make([]option.ContentOption, 0, maxAnswers+1)
	for i := range maxAnswers + 1 {
		tooManyAnswers = append(tooManyAnswers, makeAnswer(fmt.Sprintf("%d", i)))
	}
	invalidFormat, _ := option.NewContentOption(option.ContentTypeAudio, "audio")

	tests := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "success valid short question",
			params: Params{
				Text: makeText("Назовите столицу Франции"),
				Answers: []option.ContentOption{
					makeAnswer("Париж"),
					makeAnswer("париж"),
				},
			},
			err: nil,
		},
		{
			name: "success with duplicates allowed",
			params: Params{
				Text: makeText("Тест дубликатов"),
				Answers: []option.ContentOption{
					makeAnswer("Да"),
					makeAnswer("Да"),
				},
			},
			err: ErrDuplicateAnswer,
		},
		{
			name: "error no answers",
			params: Params{
				Text:    makeText("Вопрос без ответа"),
				Answers: []option.ContentOption{},
			},
			err: ErrNoAnswers,
		},
		{
			name: "error too many answers",
			params: Params{
				Text:    makeText("Слишком много вариантов"),
				Answers: tooManyAnswers,
			},
			err: ErrTooManyAnswers,
		},
		{
			name: "error invalid answer format",
			params: Params{
				Text:    makeText("Недопустимый вариант ответа"),
				Answers: []option.ContentOption{invalidFormat},
			},
			err: ErrInvalidAnswerFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, q)
				shortQ, ok := q.(*ShortQuestion)
				require.True(t, ok)

				assert.Equal(t, shortQ.Answers(), tt.params.Answers)
				assert.Equal(t, question.TypeShort, shortQ.Type())
			}
		})
	}
}

func TestShortQuestion_UpdateAnswers(t *testing.T) {
	params := Params{
		Text: "Базовый текст",
		Answers: []option.ContentOption{
			makeAnswer("Ответ 1"),
		},
	}
	q, err := New(params)
	require.NoError(t, err)

	shortQ, ok := q.(*ShortQuestion)
	require.True(t, ok)

	t.Run("success update answers", func(t *testing.T) {
		newAnswers := []option.ContentOption{
			makeAnswer("Новый 1"),
			makeAnswer("Новый 2"),
		}

		err := shortQ.UpdateAnswers(newAnswers)

		assert.NoError(t, err)
		assert.Equal(t, newAnswers, shortQ.Answers())
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldAnswers := shortQ.Answers()

		invalidAnswers := []option.ContentOption{
			makeAnswer("Норм"),
			makeAnswer("Норм"),
		}
		err := shortQ.UpdateAnswers(invalidAnswers)

		assert.ErrorIs(t, err, ErrDuplicateAnswer)
		assert.Equal(t, oldAnswers, shortQ.Answers())
	})
}

func TestHasAnswer(t *testing.T) {
	params := Params{
		Text: "Базовый текст",
		Answers: []option.ContentOption{
			makeAnswer("Ответ 1"),
		},
	}
	q, err := New(params)
	require.NoError(t, err)

	shortQ, ok := q.(*ShortQuestion)
	require.True(t, ok)

	t.Run("has answer", func(t *testing.T) {
		assert.True(t, shortQ.HasAnswer(makeAnswer("Ответ 1")))
	})

	t.Run("has answer with different case", func(t *testing.T) {
		assert.True(t, shortQ.HasAnswer(makeAnswer("ответ 1")))
		assert.True(t, shortQ.HasAnswer(makeAnswer("ОТВЕТ 1")))
	})

	t.Run("has not answer", func(t *testing.T) {
		assert.False(t, shortQ.HasAnswer(makeAnswer("Ответ 2")))
	})
}

func makeText(s string) question.QText {
	text, _ := question.NewQText(s)
	return text
}

func makeAnswer(s string) option.ContentOption {
	answer, _ := option.NewContentOption(option.ContentTypeText, s)
	return answer
}
