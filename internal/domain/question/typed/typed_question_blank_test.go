package typed

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
)

func TestNewBlankParams(t *testing.T) {
	tooManyAnswers := make([]option.ContentOption, 0, maxAnswersPerPlaceholder+1)
	for i := range maxAnswersPerPlaceholder + 1 {
		tooManyAnswers = append(tooManyAnswers, makeOption(fmt.Sprintf("%d", i)))
	}

	nonTextOpt, _ := option.NewContentOption(option.ContentTypeImage, "image")

	tests := []struct {
		name   string
		params BlankParams
		err    error
	}{
		{
			name: "empty placeholder",
			params: BlankParams{
				Placeholder: "",
				Answers: []option.ContentOption{
					makeOption("1"),
					makeOption("2"),
					makeOption("3"),
				},
			},
			err: ErrInvalidPlaceholder,
		},
		{
			name: "non-regexp placeholder",
			params: BlankParams{
				Placeholder: "{placeholder}",
				Answers: []option.ContentOption{
					makeOption("1"),
					makeOption("2"),
					makeOption("3"),
				},
			},
			err: ErrInvalidPlaceholder,
		},
		{
			name: "empty answers",
			params: BlankParams{
				Placeholder: "{{placeholder}}",
				Answers:     []option.ContentOption{},
			},
			err: ErrEmptyAnswers,
		},
		{
			name: "too many answers",
			params: BlankParams{
				Placeholder: "{{placeholder}}",
				Answers:     tooManyAnswers,
			},
			err: ErrTooManyAnswers,
		},
		{
			name: "invalid content type of answer",
			params: BlankParams{
				Placeholder: "{{placeholder}}",
				Answers:     []option.ContentOption{nonTextOpt},
			},
			err: ErrInvalidAnswerFormat,
		},
		{
			name: "valid params",
			params: BlankParams{
				Placeholder: "{{placeholder}}",
				Answers:     []option.ContentOption{makeOption("1"), makeOption("2")},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewBlank(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if err == nil {
				assert.Equal(t, p.Placeholder(), tt.params.Placeholder)
				assert.Equal(t, p.Answers(), tt.params.Answers)
			}
		})
	}
}
