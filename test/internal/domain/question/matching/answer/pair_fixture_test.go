package answer_test

import (
	"gitflic.ru/lms/internal/domain/question/matching/answer"
	"github.com/google/uuid"
)

func makeAnswerPair(promptID, matchID uuid.UUID) answer.AnswerPair {
	return answer.AnswerPair{
		PromptID: promptID,
		MatchID:  matchID,
	}
}
