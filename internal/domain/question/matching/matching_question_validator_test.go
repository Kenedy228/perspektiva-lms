package matching

import (
	"errors"
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/content"
	"github.com/google/uuid"
)

func TestValidatePairPrompt(t *testing.T) {
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
			prompt: "   ",
			err:    ErrEmptyPrompt,
		},
		{
			name:   "valid prompt",
			prompt: "prompt",
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePairPrompt(tt.prompt)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}

func TestValidatePairOption(t *testing.T) {
	tests := []struct {
		name   string
		option uuid.UUID
		err    error
	}{
		{
			name:   "nil option",
			option: uuid.Nil,
			err:    ErrNilOption,
		},
		{
			name:   "valid option",
			option: uuid.New(),
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePairOption(tt.option)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}

func TestValidatePairs(t *testing.T) {
	type fooContent struct {
		prompt string
		cType  content.ContentType
		value  string
	}

	tests := []struct {
		name       string
		content    []fooContent
		pairsCount int
		err        error
	}{
		{
			name:       "empty pairs",
			content:    []fooContent{},
			pairsCount: 0,
			err:        ErrEmptyPairs,
		},
		{
			name:       "duplicated prompts",
			content:    []fooContent{{prompt: "p1", cType: content.ContentTypeImage, value: "val"}, {prompt: "p", cType: content.ContentTypeImage, value: "val2"}, {prompt: "p", cType: content.ContentTypeImage, value: "val2"}},
			pairsCount: 3,
			err:        ErrPromptDuplicate,
		},
		{
			name:       "extra prompts",
			content:    []fooContent{{prompt: "p", cType: content.ContentTypeImage, value: "val"}, {prompt: "p1", cType: content.ContentTypeImage, value: "val2"}, {prompt: "p2", cType: content.ContentTypeImage, value: "val2"}},
			pairsCount: 2,
			err:        ErrExtraPrompts,
		},
		{
			name:       "duplicated option",
			content:    []fooContent{{prompt: "p", cType: content.ContentTypeImage, value: "val"}, {prompt: "p2", cType: content.ContentTypeImage, value: "val"}},
			pairsCount: 2,
			err:        ErrOptionDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs := make(map[string]content.RichContent, len(tt.content))

			for i := range tt.content {
				c, err := content.New(tt.content[i].cType, tt.content[i].value)
				if err != nil {
					t.Errorf("expected err nil, got %v", err)
				}

				pairs[tt.content[i].prompt] = c
			}

			err := validatePairs(pairs, tt.pairsCount)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}

func TestSizeLimitExceeded(t *testing.T) {
	tests := []struct {
		name string
		size int
		err  error
	}{
		{
			name: "size 19",
			size: 19,
			err:  nil,
		},
		{
			name: "size 20",
			size: 20,
			err:  nil,
		},
		{
			name: "size 21",
			size: 21,
			err:  ErrTooManyPairs,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs := make(map[string]content.RichContent, tt.size)

			for i := range tt.size {
				c, err := content.New(content.ContentTypeImage, fmt.Sprintf("abc%d", i))
				if err != nil {
					t.Errorf("expected err nil, got %v", err)
				}

				pairs[fmt.Sprintf("%d", i)] = c
			}

			err := validatePairs(pairs, len(pairs))

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}
