package matching

import (
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Pair struct {
	id     uuid.UUID
	prompt string
	option uuid.UUID
}

func NewPair(prompt string, optionID uuid.UUID) (Pair, error) {
	if err := validatePairPrompt(prompt); err != nil {
		return Pair{}, err
	}

	if err := validatePairOption(optionID); err != nil {
		return Pair{}, err
	}

	id, err := utils.GenerateID()
	if err != nil {
		return Pair{}, err
	}

	return Pair{
		id:     id,
		prompt: prompt,
		option: optionID,
	}, nil
}

func (p Pair) ID() uuid.UUID {
	return p.id
}

func (p Pair) Prompt() string {
	return p.prompt
}

func (p Pair) Option() uuid.UUID {
	return p.option
}
