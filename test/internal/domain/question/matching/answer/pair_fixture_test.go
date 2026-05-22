//go:build legacy
// +build legacy

package answer_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"github.com/google/uuid"
)

func makeAnswerPair(promptID, matchID uuid.UUID) answer.Pair {
	return answer.Pair{
		PromptID: promptID,
		MatchID:  matchID,
	}
}
