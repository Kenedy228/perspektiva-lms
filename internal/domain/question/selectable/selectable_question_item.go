package selectable

import (
	"gitflic.ru/lms/internal/domain/question/option"
)

type Item struct {
	content   option.ContentOption
	isCorrect bool
}

func NewItem(params ItemParams) Item {
	return Item{
		content:   params.Content,
		isCorrect: params.IsCorrect,
	}
}

func (i Item) Content() option.ContentOption {
	return i.content
}

func (i Item) IsCorrect() bool {
	return i.isCorrect
}

func (i Item) Equal(other Item) bool {
	return i.content == other.content
}
