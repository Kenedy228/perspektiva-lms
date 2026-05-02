package matching

import (
	"gitflic.ru/lms/internal/domain/grading/score"
	"gitflic.ru/lms/internal/domain/question"
	qmatching "gitflic.ru/lms/internal/domain/question/matching"
	"gitflic.ru/lms/internal/domain/question/matching/answer"
)

type Checker struct {
}

func New() Checker {
	return Checker{}
}

func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	qCast, ok := q.(*qmatching.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	pairs := qCast.Pairs()
	studentAnswers := aCast.PairsAsMap()

	if len(studentAnswers) > len(pairs) {
		s, _ := score.New(0)
		return s, nil
	}

	for i := range pairs {
		matchID, ok := studentAnswers[pairs[i].PromptID()]

		if !ok {
			s, _ := score.New(0)
			return s, nil
		}

		if ok && pairs[i].MatchID() != matchID {
			s, _ := score.New(0)
			return s, nil
		}
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question.Type) bool {
	return qType == question.TypeMatching
}
