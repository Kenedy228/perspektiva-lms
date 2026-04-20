package sequence

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
	"gitflic.ru/lms/internal/domain/question/option"
)

const minItems = 2
const maxItems = 20

type SequenceQuestion struct {
	base.Base
	items []option.ContentOption
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateItems(params.Items); err != nil {
		return nil, err
	}

	cItems := slices.Clone(params.Items)

	return &SequenceQuestion{
		Base:  base,
		items: cItems,
	}, nil
}

func (q *SequenceQuestion) Items() []option.ContentOption {
	return slices.Clone(q.items)
}

func (q *SequenceQuestion) ShuffledItems() []option.ContentOption {
	cItems := slices.Clone(q.items)
	rand.Shuffle(len(cItems), func(i, j int) {
		cItems[i], cItems[j] = cItems[j], cItems[i]
	})

	return cItems
}

func (q *SequenceQuestion) Type() question.Type {
	return question.TypeSequence
}

func (q *SequenceQuestion) UpdateItems(rawItems []option.ContentOption) error {
	if err := validateItems(rawItems); err != nil {
		return err
	}

	cItems := slices.Clone(rawItems)

	q.items = cItems
	q.Touch()
	return nil
}

func (q *SequenceQuestion) HasItem(item option.ContentOption) bool {
	return slices.Contains(q.items, item)
}
