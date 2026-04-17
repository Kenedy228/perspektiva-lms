package matching

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tooManyPairs := make([]PairParam, maxPairs+1)
	for i := 0; i <= maxPairs; i++ {
		tooManyPairs = append(tooManyPairs, makePairParam(fmt.Sprintf("prompt-%d", i), fmt.Sprintf("%d", i)))
	}

	tests := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "success",
			params: Params{
				Text: createText("Сопоставьте страны и столицы"),
				Pairs: []PairParam{
					makePairParam("Россия", "1"),
					makePairParam("Франция", "2"),
				},
			},
			err: nil,
		},
		{
			name: "empty pairs",
			params: Params{
				Text:  createText("text"),
				Pairs: []PairParam{},
			},
			err: ErrEmptyPairs,
		},
		{
			name: "not enough pairs",
			params: Params{
				Text: createText("text"),
				Pairs: []PairParam{
					makePairParam("Россия", "1"),
				},
			},
			err: ErrNotEnoughPairs,
		},
		{
			name: "too many pairs",
			params: Params{
				Text:  createText("text"),
				Pairs: tooManyPairs,
			},
			err: ErrTooManyPairs,
		},
		{
			name: "duplicate option",
			params: Params{
				Text: createText("text"),
				Pairs: []PairParam{
					makePairParam("Россия", "1"),
					makePairParam("Франция", "1"),
				},
			},
			err: ErrDuplicateOption,
		},
		{
			name: "prompt duplicate detected by count",
			params: Params{
				Text: createText("text"),
				Pairs: []PairParam{
					makePairParam("Россия", "1"),
					makePairParam("Россия", "5"),
					makePairParam("Франция", "2"),
				},
			},
			err: ErrDuplicatePrompt,
		},
		{
			name: "empty prompt in pair",
			params: Params{
				Text: "text",
				Pairs: []PairParam{
					makePairParam("", "1"),
					makePairParam("Франция", "2"),
				},
			},
			err: ErrEmptyPrompt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err != nil {
				assert.Nil(t, q)
				return
			}

			require.NotNil(t, q)

			mq, ok := q.(*MatchingQuestion)
			require.True(t, ok)

			assert.Equal(t, question.TypeMatching, mq.Type())
			assert.Len(t, mq.Pairs(), len(tt.params.Pairs))
			if len(tt.params.Pairs) != 0 {
				assert.Equal(t, mq.Pairs(), convertParamsToPairs(tt.params.Pairs))
			}
		})
	}
}

func TestMatchingQuestion_UpdatePairs(t *testing.T) {
	q, err := New(Params{
		Text: createText("Сопоставьте страны и столицы"),
		Pairs: []PairParam{
			makePairParam("Россия", "1"),
			makePairParam("Франция", "2"),
		},
	})
	require.NoError(t, err)

	mq, ok := q.(*MatchingQuestion)
	require.True(t, ok)

	t.Run("success", func(t *testing.T) {
		oldUpdatedAt := mq.UpdatedAt()
		err := mq.UpdatePairs([]PairParam{
			makePairParam("Германия", "1"),
			makePairParam("Италия", "2"),
			makePairParam("Испания", "3"),
		})

		require.NoError(t, err)
		assert.Len(t, mq.Pairs(), 3)
		assert.True(t, oldUpdatedAt.Before(mq.UpdatedAt()))
	})

	t.Run("error keeps old state", func(t *testing.T) {
		oldPairs := mq.Pairs()

		err := mq.UpdatePairs([]PairParam{
			makePairParam("Только один", "10"),
		})

		assert.ErrorIs(t, err, ErrNotEnoughPairs)
		assert.Equal(t, oldPairs, mq.Pairs())
	})
}

func TestHasPair(t *testing.T) {
	q, err := New(Params{
		Text: createText("Сопоставьте страны и столицы"),
		Pairs: []PairParam{
			makePairParam("Россия", "1"),
			makePairParam("Франция", "2"),
		},
	})
	require.NoError(t, err)

	mq, ok := q.(*MatchingQuestion)
	require.True(t, ok)

	t.Run("has pair", func(t *testing.T) {
		target := makePair("Россия", "1")

		assert.True(t, mq.HasPair(target))
	})

	t.Run("has not pair", func(t *testing.T) {
		diffPromptAndOption := makePair("Япония", "3")
		diffPrompt := makePair("Япония", "1")
		diffOption := makePair("Россия", "5")

		assert.False(t, mq.HasPair(diffPromptAndOption))
		assert.False(t, mq.HasPair(diffPrompt))
		assert.False(t, mq.HasPair(diffOption))
	})
}

func makeOption(s string) option.ContentOption {
	opt, _ := option.NewContentOption(option.ContentTypeText, s)
	return opt
}

func makePair(prompt, s string) Pair {
	opt := makeOption(s)
	pair, _ := NewPair(prompt, opt)
	return pair
}

func makePairParam(prompt, s string) PairParam {
	opt := makeOption(s)
	return PairParam{
		Prompt:        prompt,
		ContentOption: opt,
	}
}

func convertParamsToPairs(params []PairParam) []Pair {
	pairs := make([]Pair, 0, len(params))

	for i := range params {
		p := makePair(params[i].Prompt, params[i].ContentOption.Value())
		pairs = append(pairs, p)
	}

	return pairs
}

func createText(s string) question.QText {
	text, _ := question.NewQText(s)
	return text
}
