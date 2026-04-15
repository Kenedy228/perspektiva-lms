package sequence

import (
	"slices"

	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

const description = "расставьте события в хронологическом порядке"
const minItems = 2
const maxItems = 20

type SequenceQuestion struct {
	base.Base
	items []Item
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

	if err := validateItems(params.Items); err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(params.Items))
	for i := range params.Items {
		item, err := NewItem(params.Items[i])
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &SequenceQuestion{
		Base:  base,
		items: items,
	}, nil
}

func (q *SequenceQuestion) Items() []Item {
	return slices.Clone(q.items)
}

func (q *SequenceQuestion) UpdateItems(rawItems []content.RichContent) error {
	if err := validateItems(rawItems); err != nil {
		return err
	}

	items := make([]Item, 0, len(rawItems))

	for i := range rawItems {
		item, err := NewItem(rawItems[i])
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	q.items = items
	q.Touch()
	return nil
}

func (q *SequenceQuestion) Type() question.Type {
	return question.TypeSequence
}
