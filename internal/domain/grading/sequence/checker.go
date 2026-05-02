package sequence

import (
	"gitflic.ru/lms/internal/domain/grading/score"
	"gitflic.ru/lms/internal/domain/question"
	qsequence "gitflic.ru/lms/internal/domain/question/sequence"
	"gitflic.ru/lms/internal/domain/question/sequence/answer"
)

type Checker struct {
}

func New() Checker {
	return Checker{}
}

func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	qCast, ok := q.(*qsequence.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	seq := qCast.Options()
	studentAnswers := aCast.Options()

	if len(seq) != len(studentAnswers) {
		s, _ := score.New(0)
		return s, nil
	}

	for i := range seq {
		if seq[i].ID() != studentAnswers[i].OptionID {
			s, _ := score.New(0)
			return s, nil
		}
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question.Type) bool {
	return qType == question.TypeSequence
}
