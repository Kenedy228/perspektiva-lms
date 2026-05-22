package matching

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"github.com/google/uuid"
)

type Checker struct {
}

func New() Checker {
	return Checker{}
}

func (c Checker) Check(q question2.Question, a question2.Answer) (score.Score, error) {
	qCast, ok := q.(*qmatching.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	pairs := qCast.Pairs()
	studentAnswers := aCast.Pairs()

	if len(studentAnswers) > len(pairs) {
		s, _ := score.New(0)
		return s, nil
	}

	for i := range pairs {
		if !containsPair(studentAnswers, pairs[i].PromptID(), pairs[i].MatchID()) {
			s, _ := score.New(0)
			return s, nil
		}
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question2.Type) bool {
	return qType == question2.TypeMatching
}

func containsPair(pairs []answer.Pair, promptID, matchID uuid.UUID) bool {
	for i := range pairs {
		if pairs[i].PromptID.ID() == promptID && pairs[i].MatchID.ID() == matchID {
			return true
		}
	}

	return false
}
