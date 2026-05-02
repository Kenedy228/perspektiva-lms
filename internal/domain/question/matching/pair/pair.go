package pair

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/matching/pair/element"
	"github.com/google/uuid"
)

type Pair struct {
	prompt element.Element
	match  element.Element
}

func New(prompt, match content.Content) (Pair, error) {
	if err := validatePrompt(prompt); err != nil {
		return Pair{}, err
	}

	if err := validateMatch(match); err != nil {
		return Pair{}, err
	}

	promptElement, err := element.New(prompt)
	if err != nil {
		return Pair{}, err
	}

	matchElement, err := element.New(match)
	if err != nil {
		return Pair{}, err
	}

	return Pair{
		prompt: promptElement,
		match:  matchElement,
	}, nil
}

func (p Pair) Prompt() element.Element {
	return p.prompt
}

func (p Pair) Match() element.Element {
	return p.match
}

func (p Pair) PromptID() uuid.UUID {
	return p.prompt.ID()
}

func (p Pair) MatchID() uuid.UUID {
	return p.match.ID()
}
