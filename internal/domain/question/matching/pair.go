package matching

import (
	"gitflic.ru/lms/internal/domain/question"
)

type Pair struct {
	prompt  string
	content question.Content
}

func NewPair(params PairParams) (Pair, error) {
	if err := validatePrompt(params.Prompt); err != nil {
		return Pair{}, err
	}

	return Pair{
		prompt: params.Prompt,
		content: params.Content,
	}, nil
}

func (p Pair) Prompt() string {
	return p.prompt
}

func (p Pair) Content() question.Content {
	return p.content
}

func (p Pair) Equal(other Pair) bool {
	return p.prompt == other.prompt && p.content == other.content
}
