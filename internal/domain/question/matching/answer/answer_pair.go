package answer

import "github.com/google/uuid"

type AnswerPair struct {
	PromptID uuid.UUID
	MatchID  uuid.UUID
}
