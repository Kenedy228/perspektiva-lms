package selectable

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
)

type Checker struct {
}

func New() Checker {
	return Checker{}
}

func (c Checker) Check(q question2.Question, a question2.Answer) (score.Score, error) {
	qCast, ok := q.(*qselectable.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	opts := qCast.Options()
	studentAnswers := aCast.OptionIDSet()

	if len(studentAnswers) > qCast.CorrectOptionsCount() {
		s, _ := score.New(0)
		return s, nil
	}

	correct := 0
	for i := range opts {
		_, ok := studentAnswers[opts[i].ID()]

		if ok {
			if opts[i].IsCorrect() {
				correct++
			} else {
				s, _ := score.New(0)
				return s, nil
			}
		}
	}

	if correct != qCast.CorrectOptionsCount() {
		s, _ := score.New(0)
		return s, nil
	}

	s, _ := score.New(1)
	return s, nil
}

func (c Checker) Supports(qType question2.Type) bool {
	return qType == question2.TypeSelectable
}
