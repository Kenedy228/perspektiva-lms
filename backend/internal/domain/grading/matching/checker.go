package matching

import (
	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"github.com/google/uuid"
)

// Checker проверяет ответы для вопросов на сопоставление пар.
type Checker struct{}

// New создает checker для вопросов типа matching.
func New() Checker {
	return Checker{}
}

// Check возвращает 1 только когда все пары вопроса присутствуют в ответе.
func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	qCast, ok := q.(*matching.Question)
	if !ok {
		return score.Score{}, grading.ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, grading.ErrInvalidAnswerType
	}

	pairs := qCast.Pairs()
	studentAnswers := aCast.Pairs()

	if len(studentAnswers) > len(pairs) {
		return score.New(0)
	}

	for i := range pairs {
		if !containsPair(studentAnswers, pairs[i].PromptID(), pairs[i].MatchID()) {
			return score.New(0)
		}
	}

	return score.New(1)
}

// Supports сообщает, поддерживается ли тип вопроса checker-ом.
func (c Checker) Supports(qType question.Type) bool {
	return qType == question.TypeMatching
}

func containsPair(pairs []answer.Pair, promptID, matchID uuid.UUID) bool {
	for i := range pairs {
		if pairs[i].PromptID == promptID && pairs[i].MatchID == matchID {
			return true
		}
	}

	return false
}
