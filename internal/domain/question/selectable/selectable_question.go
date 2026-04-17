package selectable

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const minItems = 2
const maxItems = 20
const minCorrectItems = 1

type SelectableQuestion struct {
	base.Base
	items []Item
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateItems(params.Items); err != nil {
		return nil, err
	}

	items := mapItems(params.Items)

	return &SelectableQuestion{
		Base:  base,
		items: items,
	}, nil
}

func (q *SelectableQuestion) Items() []Item {
	return slices.Clone(q.items)
}

func (q *SelectableQuestion) ShuffledItems() []Item {
	cItems := slices.Clone(q.items)
	rand.Shuffle(len(cItems), func(i, j int) {
		cItems[i], cItems[j] = cItems[j], cItems[i]
	})

	return cItems
}

func (q *SelectableQuestion) Type() question.Type {
	return question.TypeSelectable
}

func (q *SelectableQuestion) UpdateItems(rawItems []ItemParams) error {
	if err := validateItems(rawItems); err != nil {
		return err
	}

	items := mapItems(rawItems)

	q.items = items
	q.Touch()
	return nil
}

func (q *SelectableQuestion) HasItem(item Item) bool {
	return slices.ContainsFunc(q.items, func(i Item) bool {
		return i.Equal(item)
	})
}
