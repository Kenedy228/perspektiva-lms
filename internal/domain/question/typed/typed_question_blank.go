package typed

import (
	"slices"

	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Blank struct {
	id      uuid.UUID
	mark    string
	answers []string
}

func NewBlank(mark string, answers []string) (Blank, error) {
	if err := validateBlank(mark, answers); err != nil {
		return Blank{}, err
	}

	answersCopy := slices.Clone(answers)

	id, err := utils.GenerateID()
	if err != nil {
		return Blank{}, err
	}

	return Blank{
		id:      id,
		mark:    mark,
		answers: answersCopy,
	}, nil
}

func (b Blank) ID() uuid.UUID {
	return b.id
}

func (b Blank) Mark() string {
	return b.mark
}

func (b Blank) Answers() []string {
	return slices.Clone(b.answers)
}
