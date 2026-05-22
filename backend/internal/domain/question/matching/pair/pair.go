package pair

import "github.com/google/uuid"

type Pair struct {
	prompt Prompt
	match  Match
}

func New(prompt Prompt, match Match) (Pair, error) {
	if err := validatePrompt(prompt); err != nil {
		return Pair{}, err
	}

	if err := validateMatch(match); err != nil {
		return Pair{}, err
	}

	return Pair{
		prompt: prompt,
		match:  match,
	}, nil
}

func (p Pair) Prompt() Prompt {
	return p.prompt
}

func (p Pair) Match() Match {
	return p.match
}

func (p Pair) PromptID() uuid.UUID {
	return p.prompt.ID()
}

func (p Pair) MatchID() uuid.UUID {
	return p.match.ID()
}

func (p Pair) IsIncomplete() bool {
	return p.prompt.IsIncomplete() || p.match.IsIncomplete()
}
