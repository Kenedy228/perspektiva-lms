package matching

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPair(t *testing.T) {
	baseOption, err := option.NewContentOption(option.ContentTypeText, "text")
	require.Nil(t, err)

	tests := []struct {
		name   string
		prompt string
		err    error
	}{
		{
			name:   "empty prompt",
			prompt: "",
			err:    ErrEmptyPrompt,
		},
		{
			name:   "whitespaces prompt",
			prompt: " ",
			err:    ErrEmptyPrompt,
		},
		{
			name:   "valid prompt",
			prompt: "valid",
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pair, err := NewPair(PairParams{
				Prompt:        tt.prompt,
				ContentOption: baseOption,
			})

			assert.ErrorIs(t, err, tt.err)

			if err == nil {
				assert.Equal(t, pair.Prompt(), tt.prompt)
				assert.Equal(t, pair.Option(), baseOption)
			}
		})
	}
}

func TestEqualPair(t *testing.T) {
	baseOption, err := option.NewContentOption(option.ContentTypeText, "text")
	require.Nil(t, err)
	basePair, err := NewPair(PairParams{
		Prompt:        "prompt",
		ContentOption: baseOption,
	})
	require.Nil(t, err)

	tests := []struct {
		name        string
		optionCType option.ContentType
		optionVal   string
		prompt      string
		equal       bool
	}{
		{
			name:        "same",
			optionCType: option.ContentTypeText,
			optionVal:   "text",
			prompt:      "prompt",
			equal:       true,
		},
		{
			name:        "different content option",
			optionCType: option.ContentTypeImage,
			optionVal:   "text",
			prompt:      "prompt",
			equal:       false,
		},
		{
			name:        "different prompt",
			optionCType: option.ContentTypeText,
			optionVal:   "text",
			prompt:      "another prompt",
			equal:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := option.NewContentOption(tt.optionCType, tt.optionVal)
			require.Nil(t, err)
			pair, err := NewPair(PairParams{
				Prompt:        tt.prompt,
				ContentOption: opt,
			})
			require.Nil(t, err)

			if tt.equal {
				assert.Equal(t, pair, basePair)
			} else {
				assert.NotEqual(t, pair, basePair)
			}
		})
	}
}
