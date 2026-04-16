package matching

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func contentValue(id int) content.RichContent {
	p := content.Params{
		Type: content.ContentTypeImage,
		Value: fmt.Sprintf("%d", id),
	}
	rc, _ := content.New(p)
	return rc
}

func TestNewOption(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := contentValue(1)

		opt, err := NewOption(c)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, opt.ID())
		assert.Equal(t, c, opt.Content())
	})
}

func TestNewPair(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		optionID uuid.UUID
		wantErr  error
	}{
		{
			name:     "success",
			prompt:   "Москва",
			optionID: uuid.New(),
			wantErr:  nil,
		},
		{
			name:     "empty prompt",
			prompt:   "",
			optionID: uuid.New(),
			wantErr:  ErrEmptyPrompt,
		},
		{
			name:     "whitespaces prompt",
			prompt:   "   ",
			optionID: uuid.New(),
			wantErr:  ErrEmptyPrompt,
		},
		{
			name:     "nil option",
			prompt:   "Москва",
			optionID: uuid.Nil,
			wantErr:  ErrNilOption,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pair, err := NewPair(tt.prompt, tt.optionID)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				return
			}

			assert.NotEqual(t, uuid.Nil, pair.ID())
			assert.Equal(t, tt.prompt, pair.Prompt())
			assert.Equal(t, tt.optionID, pair.Option())
		})
	}
}

func TestNew(t *testing.T) {
	tooManyPairs := make(map[string]content.RichContent, maxPairs+1)
	for i := 0; i <= maxPairs; i++ {
		tooManyPairs[fmt.Sprintf("prompt-%d", i)] = contentValue(i)
	}

	tests := []struct {
		name    string
		params  *Params
		wantErr error
	}{
		{
			name: "success",
			params: &Params{
				Text:       "Сопоставьте страны и столицы",
				PairsCount: 2,
				Pairs: map[string]content.RichContent{
					"Россия":  contentValue(1),
					"Франция": contentValue(2),
				},
			},
			wantErr: nil,
		},
		{
			name: "empty pairs",
			params: &Params{
				Text:       "text",
				PairsCount: 0,
				Pairs:      map[string]content.RichContent{},
			},
			wantErr: ErrEmptyPairs,
		},
		{
			name: "not enough pairs",
			params: &Params{
				Text:       "text",
				PairsCount: 1,
				Pairs: map[string]content.RichContent{
					"Россия": contentValue(1),
				},
			},
			wantErr: ErrNotEnoughPairs,
		},
		{
			name: "too many pairs",
			params: &Params{
				Text:       "text",
				PairsCount: maxPairs + 1,
				Pairs:      tooManyPairs,
			},
			wantErr: ErrTooManyPairs,
		},
		{
			name: "duplicate option",
			params: &Params{
				Text:       "text",
				PairsCount: 2,
				Pairs: map[string]content.RichContent{
					"Россия":  contentValue(1),
					"Франция": contentValue(1),
				},
			},
			wantErr: ErrOptionDuplicate,
		},
		{
			name: "prompt duplicate detected by count",
			params: &Params{
				Text:       "text",
				PairsCount: 3,
				Pairs: map[string]content.RichContent{
					"Россия":  contentValue(1),
					"Франция": contentValue(2),
				},
			},
			wantErr: ErrPromptDuplicate,
		},
		{
			name: "extra prompts",
			params: &Params{
				Text:       "text",
				PairsCount: 1,
				Pairs: map[string]content.RichContent{
					"Россия":  contentValue(1),
					"Франция": contentValue(2),
				},
			},
			wantErr: ErrExtraPrompts,
		},
		{
			name: "empty prompt in pair",
			params: &Params{
				Text:       "text",
				PairsCount: 2,
				Pairs: map[string]content.RichContent{
					"":        contentValue(1),
					"Франция": contentValue(2),
				},
			},
			wantErr: ErrEmptyPrompt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Nil(t, q)
				return
			}

			require.NotNil(t, q)

			mq, ok := q.(*MatchingQuestion)
			require.True(t, ok)

			assert.Equal(t, question.TypeMatching, mq.Type())
			assert.Len(t, mq.Pairs(), len(tt.params.Pairs))
			assert.Len(t, mq.Options(), len(tt.params.Pairs))

			assertPairsBoundToOptions(t, mq.Pairs(), mq.Options())
		})
	}
}

func TestMatchingQuestion_UpdatePairs(t *testing.T) {
	q, err := New(&Params{
		Text:       "Сопоставьте страны и столицы",
		PairsCount: 2,
		Pairs: map[string]content.RichContent{
			"Россия":  contentValue(1),
			"Франция": contentValue(2),
		},
	})
	require.NoError(t, err)

	mq, ok := q.(*MatchingQuestion)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		err := mq.UpdatePairs(map[string]content.RichContent{
			"Германия": contentValue(3),
			"Италия":   contentValue(4),
			"Испания":  contentValue(5),
		}, 3)

		require.NoError(t, err)
		assert.Len(t, mq.Pairs(), 3)
		assert.Len(t, mq.Options(), 3)
		assertPairsBoundToOptions(t, mq.Pairs(), mq.Options())
	})

	t.Run("error keeps old state", func(t *testing.T) {
		oldPairs := mq.Pairs()
		oldOptions := mq.Options()

		err := mq.UpdatePairs(map[string]content.RichContent{
			"Только один": contentValue(10),
		}, 1)

		assert.ErrorIs(t, err, ErrNotEnoughPairs)
		assert.Equal(t, oldPairs, mq.Pairs())
		assert.Equal(t, oldOptions, mq.Options())
	})
}

func TestMatchingQuestion_PairsEncapsulation(t *testing.T) {
	q, err := New(&Params{
		Text:       "text",
		PairsCount: 2,
		Pairs: map[string]content.RichContent{
			"Россия":  contentValue(1),
			"Франция": contentValue(2),
		},
	})
	require.NoError(t, err)

	mq := q.(*MatchingQuestion)

	got := mq.Pairs()
	require.Len(t, got, 2)

	got[0] = Pair{}

	assert.NotEqual(t, uuid.Nil, mq.Pairs()[0].ID())
}

func TestMatchingQuestion_OptionsEncapsulation(t *testing.T) {
	q, err := New(&Params{
		Text:       "text",
		PairsCount: 2,
		Pairs: map[string]content.RichContent{
			"Россия":  contentValue(1),
			"Франция": contentValue(2),
		},
	})
	require.NoError(t, err)

	mq := q.(*MatchingQuestion)

	got := mq.Options()
	require.Len(t, got, 2)

	got[0] = Option{}

	assert.NotEqual(t, uuid.Nil, mq.Options()[0].ID())
}

func assertPairsBoundToOptions(t *testing.T, pairs []Pair, options []Option) {
	t.Helper()

	optionIDs := make(map[uuid.UUID]struct{}, len(options))
	for _, opt := range options {
		assert.NotEqual(t, uuid.Nil, opt.ID())
		optionIDs[opt.ID()] = struct{}{}
	}

	for _, pair := range pairs {
		assert.NotEqual(t, uuid.Nil, pair.ID())
		assert.NotEmpty(t, pair.Prompt())

		_, ok := optionIDs[pair.Option()]
		assert.True(t, ok, "pair option id %s not found in options", pair.Option())
	}
}
