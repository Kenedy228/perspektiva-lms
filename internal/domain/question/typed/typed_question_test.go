package typed

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tooManyBlanks := make([]BlankParams, 0, maxPlaceholders+1)
	for i := range maxPlaceholders + 1 {
		placeholder := fmt.Sprintf("%d", i)
		answers := []option.ContentOption{
			makeOption(fmt.Sprintf("%d", i)),
		}
		tooManyBlanks = append(tooManyBlanks, BlankParams{
			Placeholder: placeholder,
			Answers:     answers,
		})
	}

	tests := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "success valid question",
			params: Params{
				Text: makeText("Столица России — {{moscow}}, а Франции — {{paris}}."),
				Blanks: []BlankParams{
					{
						Placeholder: "{{moscow}}",
						Answers: []option.ContentOption{
							makeOption("Москва"),
							makeOption("г. Москва"),
						},
					},
					{
						Placeholder: "{{paris}}",
						Answers: []option.ContentOption{
							makeOption("Париж"),
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "success but should be err",
			params: Params{
				Text: makeText("Столица России — {{moscow}}, а Франции — {{paris}} {{poland}}."),
				Blanks: []BlankParams{
					{
						Placeholder: "{{moscow}}",
						Answers: []option.ContentOption{
							makeOption("{{Москва}}"),
							makeOption("{{г. Москва}}"),
						},
					},
					{
						Placeholder: "{{paris}}",
						Answers: []option.ContentOption{
							makeOption("{{Париж}}"),
						},
					},
				},
			},
			err: ErrCountMismatch,
		},
		{
			name: "error no placeholders count",
			params: Params{
				Text:   makeText("Текст без пропусков"),
				Blanks: []BlankParams{},
			},
			err: ErrNotEnoughPlaceholders,
		},
		{
			name: "error too many placeholders",
			params: Params{
				Text:   makeText("{{1}}{{2}}{{3}}{{4}}{{5}}{{6}}{{7}}{{8}}{{9}}{{11}}{{12}}{{13}}{{14}}{{15}}{{16}}{{17}}{{18}}{{19}}{{21}}{{22}}{{23}}"),
				Blanks: []BlankParams{},
			},
			err: ErrTooManyPlaceholders,
		},
		{
			name: "error placeholder missing in text",
			params: Params{
				Text: makeText("Столица — {{moscow}} {{france}}."),
				Blanks: []BlankParams{
					{
						Placeholder: "{{paris}}",
						Answers: []option.ContentOption{
							makeOption("{{Париж}}"),
						},
					},
					{
						Placeholder: "{{tokyo}}",
						Answers: []option.ContentOption{
							makeOption("{{Париж}}"),
						},
					},
				},
			},
			err: ErrInvalidBlanks,
		},
		{
			name: "error duplicate mark in text",
			params: Params{
				Text: makeText("Столица — {{moscow}}. И еще раз {{moscow}}."),
				Blanks: []BlankParams{
					{
						Placeholder: "{{moscow}}",
						Answers: []option.ContentOption{
							makeOption("{{Париж}}"),
						},
					},
				},
			},
			err: ErrDuplicatePlaceholder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, q)
				typedQ, ok := q.(*TypedQuestion)
				require.True(t, ok)

				assert.Equal(t, len(typedQ.Blanks()), len(tt.params.Blanks))
				assert.Equal(t, question.TypeTyped, typedQ.Type())
			}
		})
	}
}

func TestTypedQuestion_ReplaceContent(t *testing.T) {
	params := Params{
		Text: makeText("Это {{first}} тест {{second}}."),
		Blanks: []BlankParams{
			{
				Placeholder: "{{first}}",
				Answers: []option.ContentOption{
					makeOption("{{Старый}}"),
				},
			},
			{
				Placeholder: "{{second}}",
				Answers: []option.ContentOption{
					makeOption("{{Старый}}"),
				},
			},
		},
	}
	q, err := New(params)
	require.NoError(t, err)

	typedQ := q.(*TypedQuestion)

	t.Run("success replace", func(t *testing.T) {
		err := typedQ.ReplaceContent(
			makeText("Это {{new}} контент {{first}}"),
			[]BlankParams{
				{
					Placeholder: "{{new}}",
					Answers: []option.ContentOption{
						makeOption("{{Новый}}"),
					},
				},
				{
					Placeholder: "{{first}}",
					Answers: []option.ContentOption{
						makeOption("{{Новый}}"),
					},
				},
			})

		assert.NoError(t, err)
		require.Len(t, typedQ.Blanks(), 2)
		assert.Equal(t, "{{new}}", typedQ.Blanks()[0].Placeholder())
	})

	t.Run("error invalid replace leaves state untouched", func(t *testing.T) {
		oldText := typedQ.Text()
		oldBlanks := typedQ.Blanks()

		err := typedQ.ReplaceContent(
			makeText("Текст без марков"),
			[]BlankParams{
				{
					Placeholder: "{{new}}",
					Answers: []option.ContentOption{
						makeOption("{{Новый}}"),
					},
				},
			})

		assert.ErrorIs(t, err, ErrNotEnoughPlaceholders)

		assert.Equal(t, oldText, typedQ.Text())
		assert.Equal(t, oldBlanks, typedQ.Blanks())
	})
}

func TestHasAnswerForPlaceholder(t *testing.T) {
	params := Params{
		Text: makeText("Это {{first}} тест {{second}}."),
		Blanks: []BlankParams{
			{
				Placeholder: "{{first}}",
				Answers: []option.ContentOption{
					makeOption("Старый"),
					makeOption("старый"),
					makeOption("очень старый"),
				},
			},
			{
				Placeholder: "{{second}}",
				Answers: []option.ContentOption{
					makeOption("Старый"),
					makeOption("старый"),
					makeOption("очень старый"),
				},
			},
		},
	}
	q, err := New(params)
	require.NoError(t, err)

	typedQ := q.(*TypedQuestion)

	t.Run("has answer for placeholder", func(t *testing.T) {
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("Старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("очень старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{second}}", makeOption("Старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{second}}", makeOption("старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{second}}", makeOption("очень старый")))
	})

	t.Run("has answer for mark with different cases", func(t *testing.T) {
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("старый")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("СТАРЫЙ")))
		assert.True(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("СТАРыЙ")))
	})

	t.Run("has not answer for mark", func(t *testing.T) {
		assert.False(t, typedQ.HasAnswerForPlaceholder("{{first}}", makeOption("новый")))
		assert.False(t, typedQ.HasAnswerForPlaceholder("{{third}}", makeOption("Старый")))
	})
}

func makeText(s string) question.QText {
	text, _ := question.NewQText(s)
	return text
}

func makeOption(s string) option.ContentOption {
	opt, _ := option.NewContentOption(option.ContentTypeText, s)
	return opt
}
