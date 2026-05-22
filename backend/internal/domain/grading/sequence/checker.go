package sequence

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
)

type Checker struct {
}

func New() Checker {
	return Checker{}
}

func (c Checker) Check(q question2.Question, a question2.Answer) (score.Score, error) {
	qCast, ok := q.(*qsequence.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	seq := qCast.Options()
	studentAnswers := aCast.OptionIDs()

	if len(seq) != len(studentAnswers) {
		s, _ := score.New(0)
		return s, nil
	}

	for i := range seq {
		if seq[i].ID() != studentAnswers[i].ID() {
			s, _ := score.New(0)
			return s, nil
		}
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question2.Type) bool {
	return qType == question2.TypeSequence
}
