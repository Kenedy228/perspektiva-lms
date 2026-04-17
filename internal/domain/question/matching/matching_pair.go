package matching

import (
	"gitflic.ru/lms/internal/domain/question/option"
)

type Pair struct {
	prompt string
	option option.ContentOption
}

func NewPair(params PairParams) (Pair, error) {
	if err := validatePrompt(params.Prompt); err != nil {
		return Pair{}, err
	}

	return Pair{
		prompt: params.Prompt,
		option: params.ContentOption,
	}, nil
}

func (p Pair) Prompt() string {
	return p.prompt
}

func (p Pair) Option() option.ContentOption {
	return p.option
}

func (p Pair) Equal(other Pair) bool {
	return p.prompt == other.prompt && p.option == other.option
}
