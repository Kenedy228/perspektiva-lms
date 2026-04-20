package question

import (
	"errors"
	"strings"
)

type QDescription string

const (
	QDescriptionMatching   QDescription = QDescription("сопоставьте элементы из левого списка с элементами из правого")
	QDescriptionSelectable QDescription = QDescription("выберите один или несколько правильных вариантов ответа")
	QDescriptionSequence   QDescription = QDescription("расставьте события в хронологическом порядке")
	QDescriptionShort      QDescription = QDescription("введите короткий ответ")
	QDescriptionTyped      QDescription = QDescription("заполните пропуски")
)

var (
	ErrEmptyDescription = errors.New("empty description")
)

func NewQDescription(s string) (QDescription, error) {
	if strings.TrimSpace(s) == "" {
		return QDescription(""), ErrEmptyDescription
	}

	return QDescription(s), nil
}

func (d QDescription) String() string {
	return string(d)
}

func (d QDescription) Equal(other QDescription) bool {
	return d.String() == other.String()
}
