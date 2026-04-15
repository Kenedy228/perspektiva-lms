package matching

import (
	"slices"

	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const description = "сопоставьте элементы из левого списка с элементами из правого"
const minPairs = 2
const maxPairs = 20

type MatchingQuestion struct {
	base.Base
	pairs   []Pair
	options []Option
}

func New(params *Params) (question.Question, error) {
	base, err := base.New(&base.Params{
		Text:        params.Text,
		Description: description,
		Image:       params.Image,
	})

	if err != nil {
		return nil, err
	}

	if err := validatePairs(params.Pairs, params.PairsCount); err != nil {
		return nil, err
	}

	pairs := make([]Pair, 0, len(params.Pairs))
	options := make([]Option, 0, len(params.Pairs))

	for prompt, content := range params.Pairs {
		option, err := NewOption(content)
		if err != nil {
			return nil, err
		}

		options = append(options, option)

		pair, err := NewPair(prompt, option.id)
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	return &MatchingQuestion{
		Base:    base,
		pairs:   pairs,
		options: options,
	}, nil
}

func (q *MatchingQuestion) Pairs() []Pair {
	return slices.Clone(q.pairs)
}

func (q *MatchingQuestion) UpdatePairs(rawPairs map[string]content.RichContent, pairsCount int) error {
	if err := validatePairs(rawPairs, pairsCount); err != nil {
		return err
	}

	pairs := make([]Pair, 0, len(rawPairs))
	options := make([]Option, 0, len(rawPairs))

	for prompt, content := range rawPairs {
		option, err := NewOption(content)
		if err != nil {
			return err
		}

		options = append(options, option)

		pair, err := NewPair(prompt, option.id)
		if err != nil {
			return err
		}

		pairs = append(pairs, pair)
	}

	q.pairs = pairs
	q.options = options
	q.Touch()
	return nil
}

func (q *MatchingQuestion) Options() []Option {
	return slices.Clone(q.options)
}

func (q *MatchingQuestion) Type() question.Type {
	return question.TypeMatching
}
